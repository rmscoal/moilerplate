# Moilerplate
> A monolith boilerplate for Golang backend applications with built-in strong security in mind.

> [!WARNING]<br>
> In development

## What is this project
It is a boilerplate for monolithic backend application that prioritizes security. I created this project to serves as boilerplates for my other backend applications. One of my examples project that uses this boilerplate is [SinarLog's backend](https://github.com/SinarLog/backend). It also follows Uncle Bob's Clean Architecture concepts and is inspired by some of the best clean architecture golang app out there.

## Folder Structure.
- `cmd` consists of bootstraping the app as well as starting the server.

- `config` loading the applications config by reading from .env files.

- `pkg` consists of all external/in-house packages for the application to use. Usually consists of the infrastructure or service initializations.

- `testdata` consists of mocks structs for testing.

- `internal` where all the fun begins<br>
  - `internal/domain` stores the domain of the app. I'm trying to follow Domain Driven Design as much as possible here.<br>
  - `internal/delivery` consists of the delivery methods to communicate, like the http endpoints, middlewares, and routers.<br>
  - `internal/utils` consists of application's utility functions, like primitive type manipulations.<br>
  - `internal/app` stores the application layer.<br>
    - `internal/app/usecase` consists of the application logic and orchestration of its infrastructure and service layer.<br>
    - `internal/app/repo` consists of interfaces that the infrastructure layer has to follow.<br>
    - `internal/app/service` consistes of interfaces that the service layer has to follow.<br>
  - `internal/adapter` consists of implementations to fullfil the application's infrastructure and service contracts.<br>
  - `internal/composer` acts as the manager to store all usecase, infrastructure and service layer for easier management.<br>

## How to use Moilerplate
You can start by using moilerplate in these steps:
1. Obviously, clone this project and go to the root directory of moilerplate.
2. Now here you have two options to go for the Doorkeeper's JWE<br>
    a. If you're going to use HMAC, just provide the secret key to either the `.env` file (running on your host machine) or the `docker-compose.yml` (that is if you're using docker)<br>
    b. If you're going to use aside from HMAC, like RSA for the signing, then you will need to create a `cert` folder (or any name of your choice) in the root directory. Then, for example we're going to use RSA, generate an RSA private and public key files and name it as `id_rsa` (for private) and `id_rsa.pub` (for public). The name of the file I use follows the usual namin convention. You might ask how do we generate them? Well, one way is to go to [this link](https://cryptotools.net/rsagen) for example to generate yours. Once you're done, you can input the name of the folder to `DOORKEEPER_CERT_PATH` environment variable either in `.env` or `docker-compose.yml`.<br>
    > NOTE: If you're using docker, make sure to exclude the folder you just made form the `.dockerignore` file.<br>
    > *By the way, if you guys found a better way to do is, feel free to Open PR, I'm open to solutions as long as you are using a "free" solution (not like suggestion Google Cloud Secret Manager or something 😜).*
3. Next, if you're using docker to start, I've provided a hot-reload `dev.dockerfile` for you. But if you don't want it, change the `dev.dockerfile` to `Dockerfile` in `docker-compose.yml`. Or if you're not using docker, you can go and run the app like normal: `go run .` command.


## What's already included in Moilerplate?
### Security
I believe that any system available, either available via the internet or not, should have a strong security to protect their users. While I'm building in Moilerplate, I focus a lot in security. So what security is included?
#### 1. Refresh Token reused detection
Reading an article from OAuth, [link here](https://auth0.com/blog/refresh-tokens-what-are-they-and-when-to-use-them/), it teaches the fundamentals about refresh token. Here in moilerplate, we use versioning to identifies which a refresh token is refering to. To put it short, the versioning works as below:<br>
'*Say that Alice logins in and receives AccessToken_1 and RefreshToken_1. Then at some point of time, Alice's AccessToken_1 is expired and ask for a new access token using RefreshToken_1. With that, now Alice has AccessToken_2 and RefreshToken_2. Then let's assume that there is a hacker X who stole Alice's RefreshToken_1. This hacker X then ask for a new access token with RefreshToken_1. Since we know that RefreshToken_1 has been used, the system will delete all the versioning of this token's "family". Which will require both Alice and hacker X to re-login to obtained they're access token. With this, the hacker then needs to know Alice's password for that.*'
Now, it might not be as sophisticated as other system, but this is the basic security that all application should provided.
#### 2. Rotation of the hash password
Everytime a user logs in, his/her password is being rerotated with a different random salt. Because of that the system, Moilerplate itself, does not know the salt. Also, because of this, we need an algorithm to be able to compare the hash. Moilerplate uses goroutines to look up the possible salt and uses the most recommended hashing algorithm `pbkdf2` according to my research. Furthermore, if this salt lookup takes more than **500ms**, moilerplate concludes that the given password is a mismatch. By average, Moilerplate is able to identify correct password **99.99%** of the time and it takes around **30ms** to solve correct password input. Again, once the hashed password is solved, Moilerplate generates a new salt then hash it and saves to the database.

> [!NOTE]<br>
> I also believe that there are always a room for improvement in the security applied in Moilerplate and I'm open to suggestions.

***Last but not least, I hope you enjoy using Moilerplate and always modify this boilerplate to fit your usecase!***