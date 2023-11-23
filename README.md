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
In general, Moilerplate has an in-house package named `doorkeeper` that manages password hashes and JWTs. You can customize which signing method you want to use, either symmetric or asymmetric.<br>
If you're opting for asymmetric option, you need to make a new folder that stores the private and public key. For example, if you're using RSA as your signing method, you can follow these steps:
1. Create a folder say `cert`.
2. Generate the private and public key, for example you can do via this [link](https://cryptotools.net/rsagen).
3. Paste the private key in a file named `id_rsa` inside `cert` and paste the public key in `id_rsa.pub`.
4. Then register your folder name in the `DOORKEEPER_CERT_PATH` environment variable either in `.env` or `docker-compose.yml` (for docker).

If instead you're using symmetric option you could:
1. Register the secret key in the `DOORKEEPER_SECRET_KEY` environment variable either in `.env` or `docker-compose.yml` (for docker).
2. Keep in mind that you should not fill the `DOORKEEPER_CERT_PATH` environment variable.

Next we can either start the app via docker or normal. I also provided a hot reload docker file in `dev.dockerfile`. Try it out

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