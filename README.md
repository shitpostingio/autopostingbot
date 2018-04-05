# Autoposting bot

This bot posts images and dank memes to [Shitposting](https://t.me/shitpost), automagically.

<h2>Deploying without docker</h2>

Before building autoposting bot we need to setup our system.


<h3>MySQL</h3>

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


<h3>Config.toml</h3>
Now you're ready to create the config.toml file. Simply rename config_example.toml to config.toml and edit it.


<h3>Deploy db</h3>
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

<h3>Compile the bot</h3>
Now that it's all set we're ready to build and execute the bot. Type:
```
make build
./autoposting-bot
```

