# Web Forum Powered by AWS: Backend

This repository is the backend REST API implementation for the [Web Forum Powered by AWS](https://github.com/arnavbajpai/Common-Repository-for-Web-Forum-Project) Project. The project can be deployed locally as Go REST API or deployed as an AWS lambda function. However, it relies on the Amazon API Gateway for its security so the local deployment is striclty for unit test.  

## Related Repositories
- [Common Repository](https://github.com/arnavbajpai/Common-Repository-for-Web-Forum-Project): Parent repository containing overview and architecture for the overall project.  
- [Frontend Repository](https://github.com/arnavbajpai/Web-Forum-powered-by-AWS-Frontend): Built with React.js and hosted on AWS S3.
- [Database Repository](https://github.com/arnavbajpai/Web-Forum-powered-by-AWS-Database): Built on MySQL, hosted on AWS RDS.

## The Forum REST API
The forum REST API has the following structure: 

```
├── users
│   ├── GET         # GET user using alias
│   ├── POST        # Add User
│   ├── /{userID}
│       ├── POST    # Update User
│       ├── DELETE  # Mark user as deleted
├── post	    # A post could be a topic or comment. 
│   ├── GET         # GET posts by category, tag etc.
│   ├── POST        # Add post
│   ├── /{postID}
│       ├── GET     # GET post
│       ├── POST    # Update post
│       ├── DELETE  # Delete post
├── tags
│   ├── GET         # Get tag 
│   ├── POST        # Add tag 
│   ├── /{tagID}
│       ├── DELETE  # Delete Tag
```

The full OpenAPI 3.0 Spec is [here](./spec/Web-Forum-API-Spec.json). 

## Deployment

### Local
- Step 1: Install Go: Download and install Go by following the [instructions](https://go.dev/doc/install).
- Step 2. Clone this Repository. 
- Step 3. Open your terminal and navigate to the directory containing your cloned project.
- Step 4. Run `go run cmd/server/main.go`. The REST API will start at localhost:8080. 

Note: 
1. The database setup as described in [Database Repository](https://github.com/arnavbajpai/Web-Forum-powered-by-AWS-Database)
should be completed for testing. The DB user/password should be configured in internal/database/database.go. 
2. The API is designed to be deployed as lambda and protected by API gateway, so there is no authentication in local deployment. For local testing a user is hard-coded in the backend.  

### AWS 

**Deploying Lambda Function**
- Step 1: Follow the first three steps similar to local setup. Build the go project for AWS lambda and then create a zip file by running the following command at your terminal 

```
% GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap ./cmd/server
% zip myFunction.zip bootstrap
```

- Step 2: In the search bar at the top of your AWS Console, type “Lambda” and select lambda from services. 
- Step 3: Click “Create Function”. On this page provide a name for the function, select runtime as Amazon Linux 2023 and click Create Function. 
- Step 4: In the Code tab upload the zip file created in step 1. 
- Step 5: Under the Runtime Settings set the handler as: bootstrap. 
- Step 6: Go to Test tab and test the API as per the spec [here](./spec/Web-Forum-API-Spec.json). Note that API Gateway wraps the HTTP call in an event as described [here](https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html#api-gateway-simple-proxy-for-lambda-input-format). 

Note: Please refer to the General Instructions of the Web Forum project to ensure that lambda function can communicate with the RDS deployed as part of the [Database Project](https://github.com/arnavbajpai/Web-Forum-powered-by-AWS-Database)

**Exposing through API Gateway**
- Step 1: In the search bar at the top of your AWS Console, type “API Gateway” and select the service. 
- Step 2: Click on create API. Choose REST API and Import. 
- Step 3: Under Resources, for each endpoint, under each path, edit integration and choose lambda service and the lambda function created above. Also turn the lambda proxy integration to ON. 
- Step 4: Under authorizer, choose the type as Cognito and create an authorizer. This validates the JWT token returned by Cognito after Google authentication. 
- Step 5:For each endpoint, execute a test based on the [API spec](./spec/Web-Forum-API-Spec.json) 
