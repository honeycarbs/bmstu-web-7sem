FROM 'byjg/nginx-extras'

# Set working directory
WORKDIR /backend

# Copy over static assets
COPY static static/

# Copy over entrypoint
COPY nginx-entrypoint.sh /usr/local/bin/nginx-entrypoint.sh

# Copy Nginx config template
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/nginx.conf

# Nginx init
RUN ["chmod", "+x", "/usr/local/bin/nginx-entrypoint.sh"]
ENTRYPOINT ["/usr/local/bin/nginx-entrypoint.sh"]