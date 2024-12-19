import http.server

class CorsHTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header('Access-Control-Allow-Origin', '*')
        http.server.SimpleHTTPRequestHandler.end_headers(self)

def main():
    server = http.server.HTTPServer(('localhost', 9000), CorsHTTPRequestHandler)
    server.serve_forever()

if __name__ == '__main__':
    main()