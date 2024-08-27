---
layout: post
title: git-emoji 安装和配置
subtitle:
tags: [git]
comments: true
---


## 安装

```shell
$ sudo sh -c "curl https://raw.githubusercontent.com/mrowa44/emojify/master/emojify -o /usr/local/bin/emojify && chmod +x /usr/local/bin/emojify"
```

报错：
```shell
$ hello % emojify "Hey, I just :raising_hand: you, and this is :scream: , but here's my :calling: , so :telephone_receiver: me, maybe?"
Oh my! That’s a very old version of bash you’re using, we don’t support that anymore :(
Consider upgrading it or, if you must use bash 3.2.57(1)-release, download an old version of
emojify from here: https://github.com/mrowa44/emojify/blob/old_bash_support/emojify
```

更新：
```shell
$  brew install bash
```

测试：
```shell
$ emojify "Hey, I just :raising_hand: you, and this is :scream: , but here's my :calling: , so :telephone_receiver: me, maybe?"
```


```shell
$ alias gitlog='git log --oneline --color | emojify | less -r'

$ gitlog 
```
存在的问题就是每开一个终端都要执行该操作。


## 交互式客户端安装

```shell
$ npm i -g gitmoji-cli
```

测试
```shell
$ gitmoji --help
```

查看 emoji 列表
```shell
$ gitmoji -l
```

```shell
🎨 - :art: - Improve structure / format of the code.
⚡️ - :zap: - Improve performance.
🔥 - :fire: - Remove code or files.
🐛 - :bug: - Fix a bug.
🚑️ - :ambulance: - Critical hotfix.
✨ - :sparkles: - Introduce new features.
📝 - :memo: - Add or update documentation.
🚀 - :rocket: - Deploy stuff.
💄 - :lipstick: - Add or update the UI and style files.
🎉 - :tada: - Begin a project.
✅ - :white_check_mark: - Add, update, or pass tests.
🔒️ - :lock: - Fix security or privacy issues.
🔐 - :closed_lock_with_key: - Add or update secrets.
🔖 - :bookmark: - Release / Version tags.
🚨 - :rotating_light: - Fix compiler / linter warnings.
🚧 - :construction: - Work in progress.
💚 - :green_heart: - Fix CI Build.
⬇️ - :arrow_down: - Downgrade dependencies.
⬆️ - :arrow_up: - Upgrade dependencies.
📌 - :pushpin: - Pin dependencies to specific versions.
👷 - :construction_worker: - Add or update CI build system.
📈 - :chart_with_upwards_trend: - Add or update analytics or track code.
♻️ - :recycle: - Refactor code.
➕ - :heavy_plus_sign: - Add a dependency.
➖ - :heavy_minus_sign: - Remove a dependency.
🔧 - :wrench: - Add or update configuration files.
🔨 - :hammer: - Add or update development scripts.
🌐 - :globe_with_meridians: - Internationalization and localization.
✏️ - :pencil2: - Fix typos.
💩 - :poop: - Write bad code that needs to be improved.
⏪️ - :rewind: - Revert changes.
🔀 - :twisted_rightwards_arrows: - Merge branches.
📦️ - :package: - Add or update compiled files or packages.
👽️ - :alien: - Update code due to external API changes.
🚚 - :truck: - Move or rename resources (e.g.: files, paths, routes).
📄 - :page_facing_up: - Add or update license.
💥 - :boom: - Introduce breaking changes.
🍱 - :bento: - Add or update assets.
♿️ - :wheelchair: - Improve accessibility.
💡 - :bulb: - Add or update comments in source code.
🍻 - :beers: - Write code drunkenly.
💬 - :speech_balloon: - Add or update text and literals.
🗃️ - :card_file_box: - Perform database related changes.
🔊 - :loud_sound: - Add or update logs.
🔇 - :mute: - Remove logs.
👥 - :busts_in_silhouette: - Add or update contributor(s).
🚸 - :children_crossing: - Improve user experience / usability.
🏗️ - :building_construction: - Make architectural changes.
📱 - :iphone: - Work on responsive design.
🤡 - :clown_face: - Mock things.
🥚 - :egg: - Add or update an easter egg.
🙈 - :see_no_evil: - Add or update a .gitignore file.
📸 - :camera_flash: - Add or update snapshots.
⚗️ - :alembic: - Perform experiments.
🔍️ - :mag: - Improve SEO.
🏷️ - :label: - Add or update types.
🌱 - :seedling: - Add or update seed files.
🚩 - :triangular_flag_on_post: - Add, update, or remove feature flags.
🥅 - :goal_net: - Catch errors.
💫 - :dizzy: - Add or update animations and transitions.
🗑️ - :wastebasket: - Deprecate code that needs to be cleaned up.
🛂 - :passport_control: - Work on code related to authorization, roles and permissions.
🩹 - :adhesive_bandage: - Simple fix for a non-critical issue.
🧐 - :monocle_face: - Data exploration/inspection.
⚰️ - :coffin: - Remove dead code.
🧪 - :test_tube: - Add a failing test.
👔 - :necktie: - Add or update business logic.
🩺 - :stethoscope: - Add or update healthcheck.
🧱 - :bricks: - Infrastructure related changes.
🧑‍💻 - :technologist: - Improve developer experience.
💸 - :money_with_wings: - Add sponsorships or money related infrastructure.
🧵 - :thread: - Add or update code related to multithreading or concurrency.
🦺 - :safety_vest: - Add or update code related to validation.
```