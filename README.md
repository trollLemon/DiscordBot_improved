# DiscordBot_improved
A self-hostable, scalable Discord Bot written in go, with other services written in go/python.

For info on deploying the bot on your server, see [Build and Run](#build-and-run-)

Additional information for development is in [Design and Development](#design-and-development)

## Features
The bot has the following features: 

- Image Manipulation
  - Invert colors
  - Saturation 
  - Adding text
  - Edge detection
  - Erosion and dilation
  - Quality reducer 
  - Image Shuffler (split image into `n` partitions and shuffle randomly)
  - Random filter (apply a random `c*n*n` filter to an image bounded by `[x,y]`, where `c` is the color channels, and `n` is the kernel size, and `x` and `y` are user-supplied. `n` is also user supplied. )
- Image Classification 
  - Classify an image using a HuggingFace image classification model (a ViT by default)

More features to come in the future.

## Build and Run 
The following section details how to get the bot running on your server.
### Prerequisites
The entire bot application, along with its microservices, are designed to be run within docker containers. This is the recommended approach as it avoids platform incompatibilities. You will need Docker and docker-compose installed on the machine you plan to run the bot on. You will also need git to clone the repository.

First, create a bot application on the discord developer portal. [DiscordJS has a detailed section on how to get an app created](https://discordjs.guide/preparations/setting-up-a-bot-application.html#creating-your-bot). 
Make sure to keep the bot token on hand, you will need it later.

Once you create the app on the developer portal, [invite the bot to your server](https://discordjs.guide/preparations/adding-your-bot-to-servers.html#bot-invite-links). Make sure the bot application has intents to send messages, embeds, and attachments.

Copy the guildId of the server, you will need this later along with the token. See [this guide](https://support-dev.discord.com/hc/en-us/articles/360028717192-Where-can-I-find-my-Application-Team-Server-ID) for how to get the guildId.

Now that you have set up the bot application and invited it to your server, you can get the bot code ready.

Clone the repository and cd into it with 
```bash
git clone https://github.com/trollLemon/DiscordBot_improved.git && cd DiscordBot_improved
```

Next, create a `.env` file in the project root. Populate the `.env` with the following:

```.env
DISCORD_TOKEN=<discord token from earlier>
GUILD_ID=<discord server guild id>
```

Once you populate the `.env` file. Run the following to begin building and running the bot:

```bash 
docker compose up -d 
```

A few things to note:
- if your user is not in the docker group you will need to use sudo.
- The docker images (especially the image manipulation images) will take a while to build depending on your hardware. After everything is built, the bot will register the commands to your server, which may also take some time.

### Advanced Configuration
#### TODO: Explain how to use custom endpoints if the user is hosting things on different machines.


## Design and Development
This project uses a pseudo-microservice architecture. 

The core bot application is its own service, while features such as databases, image manipulation functionality, image classification functionality, are all implemented as separate containerized microservices.
#### put bot design image here


The core bot code is as decoupled as possible to allow for easy unit/integration tests. 

The following details the specification of each service:

#### GoManip
This is the microservice that provides the image manipulation functionality. It exposes rpc-like API endpoints for each function:
- `/api/image/invert/`
- `/api/image/saturate/`
- `/api/image/edgeDetection/`
- `/api/image/morphology/`
- `/api/image/reduction/`
- `/api/image/text/`
- `/api/image/randomFilter/`
- `/api/image/shuffle/`

The service doesn't provide a way to poll for image results; the service has a configurable number of workers and used vectorized operations to provide extremely fast computations for quick responses from POSTs, even under load.
Each endpoint expects a POST call with the image contained in the request body with content type `image/png` or `image/jpeg`. Parameters for the manipulation commands are supplied via query params. See GoManip's README.md for details on each endpoint's expected query params, and returned HTTP status codes.

#### Classification
This is the service that provides the classification service. The service is split into two microservices, one for exposing URL endpoints, and one for running workers on classification jobs. The services also use a redis connection to store job requests and results. The endpoint service provides two endpoints:
- `/api/v1/images`
- `/api/v1/images/classifications/{task_id}`

The first endpoint expects a POST request, with the image supplied as a multipart form data. On successful operation, the endpoint returns a `task_id` which you use to poll the result with `GET /api/v1/images/classifications/{task_id}`. See the classification README.md for information on returned HTTP status codes.

The default classification model is a ViT trained by Google from [HuggingFace](https://huggingface.co/google/vit-base-patch16-224). The model was trained on ImageNet, so if you provide the model with an obscure image, such as a meme, it will result in some humorous predictions.

## Contributing
Contributions are welcome, just make sure to adhere with the architectural patterns and write sufficient test code for your additions. 
