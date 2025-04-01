local function post_detail()
    local method = "POST"
    local path   = "/chat"  -- no query string here
    local headers = { ["Content-Type"] = "application/json" }

    -- Provide the "query" key in JSON
    local body = '{"query":"What is your return policy"}'

    return wrk.format(method, path, headers, body)
end

request = function()
    return post_detail()
end
