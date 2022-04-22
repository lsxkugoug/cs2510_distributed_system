from rpyc import Service
from rpyc.utils.server import ThreadedServer
from torchvision import models
from PIL import Image
import torch
from torchvision import transforms
import sys

def pil_loader(path):
    # open path as file to avoid ResourceWarning (https://github.com/python-pillow/Pillow/issues/835)
    with open(path, 'rb') as f:
        img = Image.open(f)
        return img.convert('RGB')

preprocess = transforms.Compose([
        transforms.Resize(256),
        transforms.CenterCrop(224),
        transforms.ToTensor(),
        transforms.Normalize(
           mean=[0.485, 0.456, 0.406],
           std=[0.229, 0.224, 0.225]
            #mean=[0.5],
           # std=[0.5]
        )])
        
img = pil_loader("cat.png")


labels = None
with open('imagenet_classes.txt') as f:
    labels = [line.strip() for line in f.readlines()]

resnet = models.resnet101(pretrained=True)


loopTimes = 1        
class TestService(Service):

    def exposed_recognition(self) -> str:
        for _ in range(loopTimes):
            img_t = preprocess(img)
            batch_t = torch.unsqueeze(img_t, 0)
            resnet.eval()
            out = resnet(batch_t)
            _, index = torch.max(out, 1)
            percentage = torch.nn.functional.softmax(out, dim=1)[0] * 100
            return labels[index[0]]

def main():
    loopTimes = sys.argv[1:][0]
    curPot = int(sys.argv[1:][1])
    print(loopTimes, curPot)
    s = ThreadedServer(TestService, port=curPot, auto_register=False)
    s.start()  
if __name__ == '__main__':
    main()

