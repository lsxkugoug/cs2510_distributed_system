import rpyc

conn = rpyc.connect('18.215.170.25', 9999)
result = conn.root.recognition()
conn.close()
print('>>>',result)

