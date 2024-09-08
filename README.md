# GoMa'luum

## i-Ma'luum scraper reimplementation with Go

<img src="https://github.com/nrmnqdds/simplified-imaluum/assets/65181897/2ad4fedc-1018-4779-b94a-5aae6f2944a3" width=100 />

## 🚧 **In Construction** 🚧
> [!IMPORTANT]
> This project is **not** associated with the official i-Ma'luum!

> [!CAUTION]
> **Not stable yet**
>
> If not working, try do it other time, later at night or anything

Support this project!

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/gbraad)

<!--A backend REST API for my infamous [Simplified i-Ma'luum](https://github.com/nrmnqdds/simplified-imaluum). Aims to improvise the performance of the scraper as Next.js server actions didn't do well in bulk parallel fetching.-->
A Reimplementation of the infamous [Simplified i-Ma'luum](https://imaluum.nrmnqdds.com) API in Go.

## Features
- [x] **Go** implementation
- [x] **Goroutine** for improved performance
- [x] **Docker** support
- [x] **Docker Compose** support

# Local installation
```
git clone http://github.com/nrmnqdds/gomaluum-api
cd gomaluum-api
go mod download
air
```

# Using Docker
```
docker build -t gomaluum .
docker run -p 1323:1323 -d gomaluum
```

# Using Docker Compose
```
docker-compose up -d gomaluum
```