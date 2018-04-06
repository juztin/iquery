# iquery

Browser based DB query tool that doesn't require any client drivers to be installed.  
Current supported DB drivers are: DB2-iSeries, DB2-Community Edition, MSSQL, MySQL, Postgres, and SQLite.

 - [Installation](#installation)
 - [Building](#building)
 - [Notes](#notes)

---

## Installation

 1. [Install Docker](#install-docker)
 2. [Download Image](#download-image)
 3. [Configuration](#configuration)
 4. [Run](#run)


### Install Docker

Download, and run, the installer for your platform:
 - **Mac**  
   https://download.docker.com/mac/stable/Docker.dmg
 - **Windows**  
   https://download.docker.com/win/stable/InstallDocker.msi
 - **Linux**  
   Your package manager will most likely have this already.  
   You can find instructions, for most distros, [here](https://docs.docker.com/engine/installation/linux/).  
    - [Debian](https://docs.docker.com/engine/installation/linux/debian/)
    - [Ubuntu](https://docs.docker.com/engine/installation/linux/ubuntu/)
    - [RedHat](https://docs.docker.com/engine/installation/linux/rhel/)
    - [CentOS](https://docs.docker.com/engine/installation/linux/centos/)


### Download Image

Pull the image from Docker Hub.
```bash
% docker pull minty/iquery
```


### Configuration
Create a script to start a container from the above image.  
The script should hold all the databases you'd like to be able to execute queries against, along with port mapping, theme, and a placeholder query.  
  
The format of the database entries are:
 > DB_X=NAME|DB_TYPE|HOSTNAME|PORT|DATABASE|USERNAME|PASSWORD


_eg._
```bash
#!/bin/bash

docker run \
	-it \
	-e "DB_1=       Users|    db2-i |     users.db.com |   446 |       users | maverick homolka |L0v{s@pplE" \
	-e "DB_1=       Users|   db2-ce |     forms.db.com | 50000 |       forms | maverick homolka |L0v{s@pplE" \
	-e "DB_3=Credit Card |    mssql |        cc.db.com |  1433 |          cc | maverick homolka |L0v{s@pplE" \
	-e "DB_3=  Customers |    mysql | customers.db.com |  3306 |   customers | maverick homolka |L0v{s@pplE" \
	-e "DB_4=  Buildings | postgres | buildings.db.com |  5432 |   buildings | maverick homolka |L0v{s@pplE" \
	-e "DB_5=       Cars |   sqlite3|                  |       | /db/test.db |                  |" \
	-e "THEME=vs-dark" \
	-e "PLACEHOLDER=SELECT * FROM SCHEMA.TABLE;" \
	-v "$HOME/tmp/db/:/db/" \
	-p 7777:8080 \
	--rm \
	minty/iquery
```


### Run

 - Execute the script created during the _Configuration_ step
 - Open your browser, and navigate to [http://localhost:7777](http://localhost:7777)
 - `cmd`+`<enter>`, or `ctrl`+`<enter>` will execute the query

---

## Building

 1. Start Container  
    ```
    ./bin/dockit.sh
    ```
 2. Build the Binary  
    _(You'll also need to pull 3rd party libs. This can be done with ```./bin/goget.sh```)_  
    ```
    ./bin/build.sh
    ```
 3. Create Image  
    ```
    docker build -t minty/iquery .
    ```

_To debug/test you can modify the `./bin/run.sh` script with valid database settings
and execute it after step 1. This will start the service and allow you to test through [http://localhost:7777](http://localhost:7777)._

---

## Notes
Changes have been made to the package: `bitbucket.org/phiggins/db2cli` to fix a few issues, _(mainly DB2 `time` types)_  
