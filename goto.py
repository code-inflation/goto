import urllib.request
import sys
import json

if len(sys.argv) < 3:
    print("Usage goto <FROM> <TO>")
    sys.exit(0)

response = urllib.request.urlopen("http://transport.opendata.ch/v1/connections?from=" + sys.argv[1] + "&to=" + sys.argv[2] + "&fields[]=from/name&fields[]=to/name&fields[]=connections/from/platform&fields[]=connections/from/departure&fields[]=connections/to/arrival").read()
j = json.loads(response)

print("Searching connection from " + j["from"]["name"] + " to " + j["to"]["name"])
print("Departure | Platform | Arrival")

for c in j["connections"]:
    print(c["from"]["departure"] + " | " + c["from"]["platform"] + " | " + c["to"]["arrival"])