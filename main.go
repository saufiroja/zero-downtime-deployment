package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Get("/", count)

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	os := map[string]operation{
		"server": func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		},
	}

	wait := gracefulShutdown(context.Background(), 1*time.Minute, os)

	<-wait

	log.Println("exiting")
}

func count(c *fiber.Ctx) error {
	server := os.Getenv("SERVER_NAME")
	time.Sleep(10 * time.Second)
	return c.SendString("App version 3.0.0, server: " + server + "\n")
}

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		// tunggu sinyal dari os
		sign := make(chan os.Signal, 1)
		signal.Notify(sign, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

		// tunggu sinyal
		<-sign

		// set timeout untuk operasi yang belum selesai
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		// tunggu semua operasi selesai
		var wg sync.WaitGroup

		// jalankan operasi
		for i, v := range ops {
			wg.Add(1)
			valueOs := v
			indexOs := i
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", indexOs)
				if err := valueOs(ctx); err != nil {
					log.Printf("%s: clean up failed: %v", indexOs, err)
					return
				}

				log.Printf("%s was shutdown gracefully", indexOs)
			}()
		}

		wg.Wait()
		close(wait)

		log.Println("shutdown completed")
	}()

	return wait
}
