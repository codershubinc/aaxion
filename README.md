[![wakatime](https://wakatime.com/badge/user/c8cd0c53-219b-4950-8025-0e666e97e8c8/project/2c075b02-2f11-41db-bb23-9ba69ced2e40.svg)](https://wakatime.com/badge/user/c8cd0c53-219b-4950-8025-0e666e97e8c8/project/2c075b02-2f11-41db-bb23-9ba69ced2e40)

# About the aaxion

this is just a project to utilize my old laptops storage (1Tb) as the my main laptop have less storage capacity (512Gb).

# Installation is simple

- check the latest release , and get the binary.
- give it executable permission `chmod +x aaxion-linux-amd64` or for the windows `aaxion-windows-amd64.exe`
- run the binary `./aaxion-linux-amd64` or `./aaxion-windows-amd64.exe`
- follow the instructions on the screen.
- enjoy!

# How it works?

`(it is implemented for linux now , windows support coming soon)`

- on linux it needs the sudo permission or you can make it `systemctl service`.
- after you made the `systemctl service` it will run on the background.
- now it checks the folders/files in root dir like `/home/swap/*`
- it lists the files/folders from root `(excluding the hidden files starting with .)`
- for api calls it uses the port `8080` by default.
- check api docs [here](./docs/api.md)

#
