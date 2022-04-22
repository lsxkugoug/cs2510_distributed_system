import rpyc

conn = rpyc.connect('localhost', 9999)
result = conn.root.recognition()
conn.close()
print('>>>',result)

