## How To Implement Zero Downtime Deployment

Example of how to implement zero downtime deployment for a simple Golang web application using Docker and Nginx as a reverse proxy. We'll use Docker Compose to orchestrate the deployment.

### Assumptions

- You have Docker and Docker Compose installed on your machine.
- You have a Golang web application that serves HTTP requests on port 8080.

### Steps

Here's the step-by-step guide:

- Create a Golang web application
- Create a Dockerfile for the Golang application
- Create a Nginx configuration file
- Create a docker-compose.yml file
- Build and start the services

### Test Zero Downtime Deployment

Now that the containers are running, your application should be accessible at http://localhost. Nginx acts as a reverse proxy, distributing requests between goservice-1 and goservice-2.
Let's simulate a zero downtime deployment by updating the Golang application:

- Make some changes to main.go to reflect the new version of your application.
- Build the Docker image again for the updated application:

```
docker-compose build goservice-1
```

- After the new image is built, update the service with the new image:

```
docker-compose up -d --no-deps --build goservice-1
```

This will update the goservice-1 service with the new image without affecting thegoservice-2 service or causing downtime.
Repeat these steps whenever you want to deploy new changes to your Golang application. By using Docker Compose and Nginx as a reverse proxy, you can achieve zero downtime deployment for your Golang web application.
