<br>

<p align="center">
    <img src="https://github.com/user-attachments/assets/3b19f6b5-9c00-4f8e-9ce1-e5e53c83edf3" align="center" width="20%">
</p>
</br>
<div align="center">
    <i>Sensor Simulator for Smart City System Development</i>
</div>
<div align="center">
<b>This CLI provides a scalable simulation of sensors, supporting any type </br> and handling thousands of concurrent units</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/henriquemarlon/congo?style=default&logo=opensourceinitiative&logoColor=white&color=1B2D3D" alt="license">
	<img src="https://img.shields.io/github/last-commit/henriquemarlon/congo?style=default&logo=git&logoColor=white&color=60CCDD" alt="last-commit">
</p>

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running](#running)


## Overview

<div align="justify">
In the process of developing systems that involve a large number of devices, such as sensors, there is a significant challenge in recreating an environment during development that closely resembles the production environment. This is not about load testing; it's about the need for a tool that can generically and highly scalably simulate a network of sensors deployed in a smart city. Congo addresses this by providing an intuitive CLI designed for use in such development environments.
</div>

## Architecture
<br>

![image](https://github.com/user-attachments/assets/39eb5ba4-af77-4eb4-a9fc-b1d12ab9fc68)

<div align="justify">
	
> The CLI communicates with a MongoDB instance where metadata for the sensors to be simulated in the production environment is stored, including attributes such as latitude, longitude, and the intervals between the minimum and maximum values the sensor can generate and publish to the system. After allocating a virtual thread for each sensor registered in the database, the system concurrently publishes the data to a topic defined and configured via an MQTT HiveMQ broker.

> The user can also make an HTTP request to the /sensor route to register a new sensor in the database, simultaneously creating a virtual thread to simulate its behavior. One potential solution for integrating this system with other services is leveraging the HiveMQ extension for Apache Kafka integration. This allows services to consume the data in an orderly fashion through a messaging queue.
</div>

##  Getting Started

###  Prerequisites

1. [Install Docker Desktop for your operating system](https://www.docker.com/products/docker-desktop/).
2. [MQTT Broker](https://www.hivemq.com/article/step-by-step-guide-using-hivemq-cloud-starter-iot/)
3. [MongoDB Instance](https://www.mongodb.com/basics/clusters/mongodb-cluster-setup)

> [!NOTE]
> For a development environment, you can use the local infrastructure provided in this repository, which includes:
>
> - **MQTT HiveMQ broker** with the Apache Kafka extension enabled.  
> - **MongoDB instance** that will be populated with the data from the provided file.  
> - Infrastructure for **Apache Kafka**.
> 
> To run this, simply clone this repository and execute the following command:
>
> ```sh
> make infra
> ```

###  Running

> [!WARNING]
> Before running the command below, ensure you have created a `Config.toml` file and set the **environment variables** correctly. **If you are using the local infrastructure referenced earlier**, use this [file](https://github.com/henriquemarlon/congo/blob/main/Config.toml) as a reference. Below is the structure of the content that should be included in the file:
>
> ```toml
> [MONGO]
> DB_URL=
> DB_NAME=
> COLLECTION_NAME=
> 
> [HIVEMQ]
> SERVER_URL=
> USERNAME=
> PASSWORD=
> ```

Run the CLI using the [Docker image](https://github.com/henriquemarlon/congo/pkgs/container/congo/330015879?tag=latest) provided for distributing its binary. Use the following command:

```sh
docker run --rm \
	-v $(pwd):/app -w /app \
	ghcr.io/henriquemarlon/congo:latest \
	--config ./Config.toml --verbose
```
