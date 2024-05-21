# aya
<img src="https://github.com/ericmariot/aya-cli/assets/29050845/7b0df4bf-5e5c-4929-b8e9-fc8b2b3b730e" alt="enhanced_image" width="200"/>

`aya` is a personal CLI in development!

## Installation
### Using Pre-built Binaries

You can download the pre-built binaries from the [releases page](https://github.com/ericmariot/aya/releases) and add the binary to your PATH.

1. Download the binary for your operating system from the releases page.
2. Move the binary to a directory included in your PATH. For example, for Linux:

   ```bash
   mv aya-linux /usr/local/bin/aya
   chmod +x /usr/local/bin/aya

## Usage
Here are some basic usage examples:

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

## Config File
-- TODO


#### License
[MIT](LICENSE)
