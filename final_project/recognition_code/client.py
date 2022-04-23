import rpyc

conn = rpyc.connect('172.18.0.1', 9999)
result = conn.root.recognition()
conn.close()
print('>>>',result)

