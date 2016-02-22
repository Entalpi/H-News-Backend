require 'rubygems'
require 'bundler/setup'
require 'securerandom'

Bundler.require(:default)

include RubyHackernews

@@users = Hash.new

helpers do

  def user
    @user ||= @@users[params[:apikey]] || halt(401)
  end

end

# Logs the user in
post '/v1/login' do
  content_type :json

  password = params[:password]
  username = params[:username]

  user = User.new(username)
  success = user.login(password)

  if success
    # Generate apikey
    apikey = SecureRandom.urlsafe_base64
    @@users[apikey] = user
    { 'apikey' => apikey }.to_json
  else
    status 404
  end
end

# Logouts a logged in user
post '/v1/login/logout' do
  apikey = params[:apikey]
  if @@user[apikey]
    @@user.delete(apikey)
  else
    status 404
  end
end

# Upvotes a post with the specified ID
post '/v1/login/entry/upvote' do
  content_type :json

  if user
    id = params[:id]
    entry = Entry.find(id)
    if entry
      entry.upvote
    else
      status 404 # Can not find Entry
    end
  end
end

# Upvotes a comment with the specified ID
post '/v1/login/comment/upvote' do
  content_type :json

  if user
    id = params[:id]
    comment = Comment.find(id)

    if comment
      comment.upvote
    else
      status 404 # Can not find Comment
    end
  end
end

# Write and post a comment on the entry with the specified ID
# Cannot write zero length comments
post '/v1/login/entry/comment' do
  content_type :json

  comment = params[:comment]
  if comment && comment.length > 0 # Valid comment?

    if user
      id = params[:id]
      entry = Entry.find(id)

      if entry
        entry.write_comment(comment)
      else
        status 404 # Can not find Entry
      end
    end
  end
end

# Write and post a reply on the comment with the specified ID
# Cannot write zero length comments
post '/v1/login/comment/reply' do
  content_type :json

  reply = params[:reply]
  if reply && reply.length > 0 # Valid reply?

    if user
      id = params[:id]
      comment = Comment.find(id)

      if comment
        comment.reply(reply)
      else
        status 404 # Can not find Comment
      end
    end
  end
end
