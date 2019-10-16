# Clickjacking PoC Tool

A basic tool to generate clickjacking proof of concepts based on a given URL.

Basic usage:

`clickjacking-poc -u https://example.com`

The tool can also be used to open the PoC up in the browser:

`Clickjacking-poc -b chromium-browser -u https://example.com`

Additional options exist for styling the template as well as supressing stdout and file output.

To see all available options and a more verbose summary of the usage of the tool run:

`clickjacking-poc -h`

## Configuration file

The tool uses the viper golang library for loading options rather than using the command line.

An example config (in JSON) format can be seen in the example-configs/ directory.

For more information on viper and the types of formats supported see https://github.com/spf13/viper

## Contributions

Accepting PR's on dev only!
