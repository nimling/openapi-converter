package converter

const locationTemplate = `{{range .Summaries}}# Summary: {{.}}
{{end}}{{range .Descriptions}}# Description: {{.}}{{end}}
location {{.Path}} {
{{if gt (len .Methods) 0}}    limit_except {{.AllowMethods}} {
        deny all;
    }
{{end}}{{if gt (len .GlobalClaims) 0}}    # Global claims required for all methods
    set_by_lua_block $missing_claims {
        local claims = ngx.var.user_claims or ""
        local missing = 0{{range .GlobalClaims}}
        if not string.match(claims, "{{.}}") then
            missing = 1
        end{{end}}
        return missing
    }
    if ($missing_claims) {
        return 403 "Missing required claims";
    }
{{end}}{{range $method, $claims := .MethodClaims}}    # Claims specific to {{$method}}
    set_by_lua_block $method_missing_claims {
        local claims = ngx.var.user_claims or ""
        if ngx.var.request_method ~= "{{$method}}" then
            return 0
        end{{range $claims}}
        if not string.match(claims, "{{.}}") then
            return 1
        end{{end}}
        return 0
    }
    if ($method_missing_claims) {
        return 403 "Missing required claims for {{$method}}";
    }
{{end}}
    rewrite ^{{.Prefix}}/(.*) /$1 break;
    proxy_pass {{.ServerURL}};

    # Basic proxy headers
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Original-URI $request_uri;

    # Azure Function specific headers
    proxy_set_header X-User-Token $http_authorization;
    proxy_set_header X-CSRF-TOKEN $http_x_csrf_token;
    proxy_set_header X-Client-IP $http_x_client_ip;

    # Timeouts
    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    proxy_read_timeout 60s;

    # Buffer settings
    proxy_buffering on;
    proxy_buffer_size 16k;
    proxy_buffers 8 16k;
    proxy_busy_buffers_size 32k;

    # HTTP/1.1 support
    proxy_http_version 1.1;
    proxy_set_header Connection "";

    # Error handling
    proxy_intercept_errors on;
    proxy_next_upstream error timeout http_500 http_502 http_503 http_504;

    # Security headers
    add_header Cache-Control "private, no-cache, no-store, must-revalidate";
    add_header Pragma no-cache;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Remove Server header
    proxy_hide_header Server;
    proxy_hide_header X-Powered-By;

    # CORS headers (if needed)
    # add_header Access-Control-Allow-Origin "*";
    # add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS";
    # add_header Access-Control-Allow-Headers "Authorization, Content-Type";
    # add_header Access-Control-Allow-Credentials "true";
}`

const oaSpecTemplate = `---
aside: false
outline: false
title: {{ .Title }}
---

<script setup lang="ts">
import { useData } from 'vitepress';
import spec from './{{ .FilePrefix }}spec.json';

const { isDark } = useData()
</script>

<OASpec :spec="spec" :isDark="isDark" />`

const oaPathsTemplate = `import { usePaths } from 'vitepress-openapi'
import spec from './{{ .FilePrefix }}spec.json' assert { type: 'json' }

export default {
    paths() {
        return usePaths({ spec })
            .getTags()
            .map(({ name }) => {
                return {
                    params: {
                        tag: name,
                        pageTitle: name
                    },
                }
            })
    },
}`

const oaTagsTemplate = `---
aside: false
outline: false
title: {{ .Title }}
---

<script setup lang="ts">
import { useRoute, useData } from 'vitepress';
import spec from './{{ .FilePrefix }}spec.json';

const route = useRoute();
const { isDark } = useData();

const tag = route.data.params.tag
</script>

<OASpec :spec="spec" :tags="[tag]" :isDark="isDark" hide-info hide-servers hide-paths-summary />`

const oaIntroductionTemplate = `---
layout: doc
title: {{ .Title }}
---

<script setup lang="ts">
import spec from './{{ .FilePrefix }}spec.json';
</script>

<OAIntroduction :spec="spec" />`
