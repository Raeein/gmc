
# gmc - give me my course

gmc is a full stack app written in go and vanilla javascript that sends you an email when there is a spot available in a course you want to take.
During course selection windows you can find the website running at https://gmc.raeeinbagheri.com

# Features

- Sends you an email when there is a spot so you dont have to email greg each semester
- CI/CD pipeline with github actions - deployment to docker hub, google cloud and binaries to github releases
- Firestore for storing contact information and emails
- MongoDB for logging as this is a serverless app

# Run the project

- Download a binary compatible with your machine from the [releases page](https://github.com/Raeein/gmc/releases)

- From Dockerhub
```bash
docker pull raeein/gmc:main-528c8eb
docker run -p 8080:8080 --name gmc --rm \
      -e SMTP_HOST=XXXX \
      -e SMTP_PORT=XXXX \
      -e SMTP_FROM=XXXX \
      -e SMTP_PASSWORD=XXXX \ 
      -e MONGO_USERNAME=XXXX \
      -e MONGO_PASSWORD=XXXX \
      -e MONGO_DATABASE=XXXX \
      -e MONGO_COLLECTION=XXXX \
      -e FIRESTORE_PROJECTID=XXXX \
      -e FIRESTORE_SECTIONS_COLLECTIONID=XXXX \
      -e FIRESTORE_USERSCOLLECTIONID=XXXX \
      -e FIRESTORE_CREDENTIALSPATH=XXXX \
      raeein/gmc:main-528c8eb
      /
```

- From source
  - You need to have `go` installed on your machine - version 1.19 recommended
  - Place your config.yaml in the root of the project
  - Have your firestore credential file
  - Fll in the config.yaml with your keys and firestore credintial path
  
```bash
git clone https://github.com/Raeein/gmc.git
cd gmc
go mod download
go run cmd/gmc/main.go
```
OR

- From root of the project 

```bash
docoker build -t gmc .  
docker run -p 8080:8080 --name gmc --rm \
      -e SMTP_HOST=XXXX \
      -e SMTP_PORT=XXXX \
      -e SMTP_FROM=XXXX \
      -e SMTP_PASSWORD=XXXX \ 
      -e MONGO_USERNAME=XXXX \
      -e MONGO_PASSWORD=XXXX \
      -e MONGO_DATABASE=XXXX \
      -e MONGO_COLLECTION=XXXX \
      -e FIRESTORE_PROJECTID=XXXX \
      -e FIRESTORE_SECTIONS_COLLECTIONID=XXXX \
      -e FIRESTORE_USERSCOLLECTIONID=XXXX \
      -e FIRESTORE_CREDENTIALSPATH=XXXX \
      gmc
      /
```

Inspiration: https://github.com/jacobmichels/Course-Sense-Go