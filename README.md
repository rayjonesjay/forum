# forum

Forum project consist in creating a web forum that allows:

1. users to communicate
2. associating categories to posts
3. liking 👍 and disliking 👎 posts and comments
4. filtering posts

## SqLite

Data (posts,users,comments..) will be stored using a SQLite.

- At Least one SELECT, ONE CREATE AND ONE INSERT queries should be used.


## Authentication

The client(forum user) should be able to `register` as a new user on the forum, by 
providing their credentials (username and password).

A `loggin session` is created in order for user to access the forum and be able to 
add posts and comments.

A user who is not logged in should not be able to add posts or comments.

`Cookies` should be used to allow each user to only have one opened session.
Each of this sessions must contain an expiration date.

It is up to you the developer to decide how long cookies should stay alive.

The use of `UUID` Universally Unique Identifier - is a 128-bit or (128/8) == 16-byte number that uniquely identifies objects in computer systems.


### Authentication Flow

1. MUST ASK FOR EMAIL
    - when email is already taken return an error response
2. MUST ASK FOR USERNAME
3. MUST ASK FOR PASSWORD
    - the password must be encrypted when stored (THIS IS A BONUS)
    - verify password

The forum backend must be able to check if the email provided is present in the database and if all credentials are correct.

## Communication

In order for users to communicate between each other, they will have to create a post and comments.

- Only registered users will be able to create posts and comments
- When registered users are creating a post they can associate one or more categories to it.
- The posts and comments should be visible to all users (registered or not).
- Non-registerd users will only be able to see posts and comments

## Likes and Dislikes
Only registered users will be able to like or dislike posts and comments
The number of likes and dislikes should be visible by all users (registered or not)

## Filter

You need to implement a filter mechanism, that will allow users to filter the displayed posts by:
- categories
- created posts
- liked posts

Subform can be used to do the filtering.

## Docker

The project will run in a docker container.
