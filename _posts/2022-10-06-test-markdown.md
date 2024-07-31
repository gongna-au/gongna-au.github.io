---
layout: post
title: What is AppImage? And how to install and use it under ubuntu? 
subtitle: AppImage 是什么？以及如何在ubuntu下面安装与使用？ 
tags: [ appImage]
---
# Linux Installation Instructions
## AppImage
layout: post
title: AppImage 是什么？以及如何在ubuntu下面安装与使用？
subtitle: Composing Methods

tags: [ linux]
---

# Linux Installation Instructions

## AppImage

[Download the AppImage](https://github.com/marktext/marktext/releases/latest) 

1. `chmod +x appName.AppImage`
2. `./appName.AppImage`
3. Now you can execute app.

### Installation

You cannot really install an AppImage. It's a file which can run directly after getting executable permission. To integrate it into desktop environment, you can either create desktop entry manually **or** use [AppImageLauncher](https://github.com/TheAssassin/AppImageLauncher).

#### Desktop file creation

See [how to create  desktop file in ubuntu ].https://www.maketecheasier.com/create-desktop-file-linux/



#### AppImageLauncher integration

You can integrate the AppImage into the system via [AppImageLauncher](https://github.com/TheAssassin/AppImageLauncher). It will handle the desktop entry automatically.

### Uninstallation

1. Delete AppImage file.
2. Delete your desktop file if exists.
3. Delete your user settings: `~/.config/appName`

### Custom launch script

1. Save AppImage somewhere. Let's say `~/bin/appname.AppImage`

2. `chmod +x ~/bin/appname.AppImage`

3. Create a launch script:

   ```sh
   #!/bin/bash
   DESKTOPINTEGRATION=0 ~/bin/appname.AppImage
   ```

