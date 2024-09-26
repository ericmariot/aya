# aya
<img src="https://github.com/ericmariot/aya-cli/assets/29050845/7b0df4bf-5e5c-4929-b8e9-fc8b2b3b730e" alt="enhanced_image" width="200"/>

`aya` is a personal CLI in development!

## Installation
To install `aya`, run the following command:

```sh
go install github.com/ericmariot/aya@latest
```

## Usage
Some basic usage examples:

```
$ aya
```
<img width="206" alt="image" src="https://github.com/ericmariot/aya/assets/29050845/b1584e76-d92c-4f57-82d8-6f35720bd18e">


#### Get Weather Information for your current location.
To get the current weather information based on your IP address, simply run:

```
$ aya weather
🌎 Getting coordinates for Criciuma
🌤️  Getting weather
America/Sao_Paulo TZ
Last update: 20:15
Current: 13.9°C
```

#### To get the current weather information for a specified city, use the weather command:
```sh
$ aya weather [city_name]
```
For example:

```sh
$ aya weather san-francisco
```
 
#### Use the Graph Flag
To display the weather information in a graphical format, use the --graph flag:

```sh
$ aya weather san-francisco --graph
```
<img width="843" alt="image" src="https://github.com/ericmariot/aya/assets/29050845/532d8bf4-ecd7-4cfd-ab99-641aab788365">

## Config File
A config file is saved at:
```
homeDir := os.UserHomeDir()
configFilePath = filepath.Join(homeDir, ".aya.json")
```

It is used to save your current location from the IP Address, and also caches cities that are used in the `weather` command!
You can clear the config file with `$ aya clearConfig`

## 🤝 Contributing

### Clone the repo

```bash
git clone https://git@github.com:ericmariot/aya.git
cd aya
```

### Build the project

```bash
go build
```

### Run the project

```bash
./aya
```

### Run the tests

```bash
go test ./...
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

#### License
[MIT](LICENSE)
