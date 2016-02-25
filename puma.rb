require 'rack-timeout'

# use Rack::Timeout          # Call as early as possible so rack-timeout runs before all other middleware.
# Rack::Timeout.timeout = 10 # Recommended. If omitted, defaults to 15 seconds.

workers Integer(ENV['WEB_CONCURRENCY'] || 2)
threads_count = Integer(ENV['MAX_THREADS'] || 5)
threads threads_count, threads_count

preload_app!

rackup      DefaultRackup
port        ENV['PORT']     || 3000
environment ENV['RACK_ENV'] || 'development'

on_worker_boot do
  Pumatra.new
end
