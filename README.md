# Set Playlist Builder

This projects aims to automate the boring task of exporting a setlist from 
[SetList.fm](https://www.setlist.fm) into [Spotify](https://www.spotify.com)

## How to run

After registering the application on Spotify, you'll get a client ID and secret key for your application. An easy way to provide this data to your application is to set the SPOTIFY_ID and SPOTIFY_SECRET environment variables. If you choose not to use environment variables, you can provide this data manually.

A config file named `setplaylistbuilder` is also needed in the same directory as the binary, with the following parameters. Json and YAML formats are accepted.

```
SetlistFmAPIKey:
SpotifyOauthTokenFile: ./spotify-oauth-token.json
```