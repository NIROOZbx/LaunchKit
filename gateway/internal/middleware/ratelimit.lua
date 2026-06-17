-- KEYS[1] = current window key
-- KEYS[2] = previous window key
-- ARGV[1] = weight (float 0..1) — how much of the previous window still counts
-- ARGV[2] = limit (max requests allowed per window)
-- ARGV[3] = ttl seconds to set on the current window key

local current  = tonumber(redis.call('GET', KEYS[1]) or '0')
local previous = tonumber(redis.call('GET', KEYS[2]) or '0')
local weight   = tonumber(ARGV[1])
local limit    = tonumber(ARGV[2])
local ttl      = tonumber(ARGV[3])

local estimated = (previous * weight) + current

if estimated >= limit then
    return 0
end

local newVal = redis.call('INCR', KEYS[1])
if newVal == 1 then
    redis.call('EXPIRE', KEYS[1], ttl)
end

return 1