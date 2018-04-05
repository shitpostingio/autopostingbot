# Autoposting bot

This bot posts images and dank memes to [Shitposting](https://t.me/shitpost), automagically.

## Deploying without Docker

Before building autoposting bot we need to setup our system.


### MySQL

Install MySQL, for example with apt, by typing in a shell. 

```
sudo apt install mysql-server
``` 

Now we'll create a database and a user. Starts by opening MySQL shell with

```
mysql -u root -p
```
and type the followings, replace databasename with everything you want:

```
create database databasename;
```

Choose username and password

```
GRANT ALL ON databasename.* TO username@localhost IDENTIFIED BY 'password';
```

Setup `utf8mb4` globally (thanks https://mathiasbynens.be/notes/mysql-utf8mb4) by editing `my.cnf` and adding - or modifying - all the relevant lines:

```
[client]
default-character-set = utf8mb4

[mysql]
default-character-set = utf8mb4

[mysqld]
character-set-client-handshake = FALSE
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
```

### Config.toml

Now you're ready to create the config.toml file. Simply rename config_example.toml to config.toml and edit it.


### Deploy db

Now that mysql is installed and the config file created we're ready to deploy the database. Navigate to database/cmd/autoposting-add-user, build with 

```
go build
```

and move the compiled binary file to the main folder of the project; do the same things in database/cmd/autoposting-deploy-db.
Move to the home of the project and execute autoposting-deploy-database by typing 

```
./autoposting-deploy-database
```

The bot will reply only to the users who knows, so to do this it needs to now your id (get it by sending a message to @rawdatabot). Execute 

```
./autoposting-add-user -userid yourid
```

to add an id to the known ones.

### Compile the bot

Now that it's all set we're ready to build and execute the bot. Type:

```
make build
./autoposting-bot
```

