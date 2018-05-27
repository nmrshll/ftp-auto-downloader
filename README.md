# ftp-auto-downloader
Auto download files from an FTP server (and delete them on the server). Tell it which files to downloads and where to put them through a really simple configuration format:

```
{
    ftpServer: "192.168.1.41:2121"
    jobs:
    [
        {
            SrcPath: /storage/emulated/0/DCIM/Camera
            DestPath: /home/me/Photos/a_trier
            AuthorizedExtensions: "jpg,gif,mp4,jpg:large,png,png:large"
        }
        {
            SrcPath: "/storage/emulated/0/WhatsApp/Media/WhatsApp Images"
            DestPath: /home/me/Photos/Whatsapp
            AuthorizedExtensions: "jpg,gif,mp4,jpg:large,png,png:large"
        }
        {
            SrcPath: "/storage/emulated/0/WhatsApp/Media/WhatsApp Video"
            DestPath: /home/me/Photos/Whatsapp
            AuthorizedExtensions: "jpg,gif,mp4,jpg:large,png,png:large"
        }
        {
            SrcPath: /storage/emulated/0/Tumblr/
            DestPath: /home/me/upload/
            AuthorizedExtensions: "jpg,gif,mp4,jpg:large,png,png:large"
            DestFolderSize: 2000
            NbFilesToLeave: 300
        }
        {
            SrcPath: /storage/emulated/0/Download/
            DestPath: /home/me/upload/
            AuthorizedExtensions: "jpg,gif,mp4,jpg:large,png,png:large"
            DestFolderSize: 2000
        }
    ]
}
```

I wrote that code to make it easier to download files from my android phone (which runs the ftp server through an app like [Software Data Cable](https://play.google.com/store/apps/details?id=com.damiapp.softdatacable&hl=en)).
Note that it deletes the files on the server after downloading them. I might make that an option later.

## Installation
### Pre-compiled binaries
Coming soon ! (or, if you've already used github releases, help is appreciated !)

### From source
#### Dependencies

 - Go 1.7+

#### Commands to run

```
go install github.com/n-marshall/ftp-auto-downloader
```

## Usage
First configure it by modifying the `config.hjson` file.

Then run `ftp-auto-downloader` in the directory where your `config.hjson` file is located.
