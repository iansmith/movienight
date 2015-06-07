### Required for launch

* Sign up page
	* So you don't have to do create_acct for every person, maybe filter on email suffix?
	* No need for oauth2 fanciness

* Delete
	* If you nominated it, you can delete it
	* On delete, up or downvotes need to deleted (cascade in DB?)
	* If you are a manager or admin, you can delete anything
	* Should probably be a control on the main page, perhaps left of the title

* Upvote and Downvote
	* The tables are in the DB for this
	* Assumption: 1 up and down vote per user at a time, thus delete movie 

* Log out
	* Pieces of this are already in place (can use same machinery as log in)

### Nice to have

* "Bringing the Pizza"
	* To drive consensus
	* Top of the main page should show the list of anybody who is the only one to up or downvote a movie

* RSVP + "showing day" (only one showing day at a time)
	* Set a "showing day" that says what movie will be shown
	* Automatically delete the move from the list when the date is in the past 
		* Super-fancy: Have a "already seen" page
	* Send email with a link that allows people to RSVP? 

* Password reset cruft
	* The scaffolding is in pwd_auth.go for this already
	* DB tables are not there though

* Oauth2 fanciness, signing with google or github
	*Probably requires seven5 changes

### Misc

* Hashed passwords 
	* https://github.com/jpillora/hashedpassword
	* https://github.com/elithrar/simple-scrypt
	* https://godoc.org/golang.org/x/crypto/bcrypt
	* Do it in s5?

* Fix the seven5 buildpack on Heroku to do the build of the client side assets and code
	* Remove all the commited "outputs" from the github repo (source)

* Ability to set a picture for your user profile
	* It's nice to convert it to a round image using the draw library of go
	* Maybe do it through github login somehow?
	* Hook to the nomination field when showing the main page
