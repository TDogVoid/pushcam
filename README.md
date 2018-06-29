# pushcam

Simple go script that is used with [motionEyeOS](https://github.com/ccrisan/motioneyeos) to send push notifications using [pushbullet](https://www.pushbullet.com/).

I didn't make a config file so you have to edit the code directly, maybe something in the future.

## Usage
Under the main func there are a couple of variables that need to be set.  You need your access toke which you can get under Account Settings on pushbullet.

set your link you want to send with push bullet (ie. the browser address of your camera)

### Build

* `GOOS=linux GOARCH=arm GOARM=5 go build pushCam.go`
* upload to \data\
* make excutable
    * `chmod +x pushCam`


### Under motioneye web interface
1. Under General Settings turn on Advanced Settings
1. Under Motion Notifications Turn on run a command
1. Enter `\data\pushCam`
1. Optional Argument you can pass is the folder for it to watch such as Camera1 or Camera2. This is usefull if you have more than one camera and want to watch for images in those cameras, otherwise it will send the first new image it finds on any camera.
    * ie `\data\pushCam Camera1`






## Warning
Pushbullet limits its api usage (at least with the free version) so if the camera goes off a lot, it will probably block it at some point.

It is also setup to send a picture in a second message, since I notice motioneyeos doesn't make the jpeg until after its done recording.  So this may not work if your only doing recording.

Sorry for the mess of coding, its something I put together quickly and haven't got around to cleaning up.