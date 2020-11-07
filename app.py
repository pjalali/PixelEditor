from flask import Flask,flash, render_template, request, redirect, url_for
import base64
import os
from werkzeug.utils import secure_filename
from PIL import Image

UPLOAD_FOLDER = '.'

app = Flask(__name__, template_folder='templates')
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER

@app.route('/')
@app.route('/index')
def index():
    return render_template('index.html')

# https://stackoverflow.com/a/29988302
@app.route('/capture', methods = ['POST'])
def capture():
    imgDataB64 = request.form['imgData']
    width = request.form['wData']
    height = request.form['hData']
    # Remove header from base64 image
    imgDataPNG = base64.b64decode(imgDataB64[22:])
    # Write image to file to comminicate with PixelEditorEngine
    with open("original.png","wb") as fo:
         fo.write(imgDataPNG)
    return render_template('edited.html', 
                            data=imgDataB64,
                            r=0, 
                            g=0, 
                            b=0,
                            c=0,
                            w=width,
                            h=height,
                            th=1,
                            hue=0,
                            sat=0,
                            light=0,
                            trgb="",
                            thsl="",
                            tt="")

@app.route('/edit', methods = ['POST'])
def edit():
    # os.system('rm output.png')
    data = request.form
    command = "./PixelEditorEngine original.png" + " " + \
              data['rOffset'] + " " + \
              data['gOffset'] + " " + \
              data['bOffset'] + " " + \
              data['contrast'] + " " + \
              data['hue'] + " " + \
              data['sat'] + " " + \
              data['light'] + " " +\
              data['threads']
    
    timesToExec = os.popen(command).read().split(' ')
    if int(data['threads']) > 1:
        timeRGB = "Parallel RGB and contrast modification using " + data['threads'] + " threads took " + timesToExec[0] + "."
    else:
        timeRGB = "Serial RGB and contrast modification took " + timesToExec[0] + "."
    timeHSL = "Serial HSL modifications took " + timesToExec[1] + "."
    timeTotal = "Total time took " + timesToExec[2] + "."
    
    with open("output.png", "rb") as img_file:
        b64Image = base64.b64encode(img_file.read())
    b64Image = "data:image/png;base64," + str(b64Image)[2:-3]

    return render_template('edited.html', 
                            data=b64Image, 
                            r=data['rOffset'], 
                            g=data['gOffset'], 
                            b=data['bOffset'],
                            c=data['contrast'],
                            w=data['widthVal'],
                            h=data['heightVal'],
                            th=data['threads'],
                            hue=data['hue'],
                            sat=data['sat'],
                            light=data['light'],
                            trgb=timeRGB,
                            thsl=timeHSL,
                            tt=timeTotal)

@app.route('/reset')
def reset():
    filename = secure_filename("original.png")
    with open(filename, "rb") as img_file:
        b64Image = base64.b64encode(img_file.read())
    b64Image = "data:image/png;base64," + str(b64Image)[2:-3]
    im = Image.open(filename)
    w, h = im.size
    return render_template('edited.html', 
                            data=b64Image,
                            r=0, 
                            g=0, 
                            b=0,
                            c=0,
                            w=w,
                            h=h,
                            th=1,
                            hue=0,
                            sat=0,
                            light=0,
                            trgb="",
                            thsl="",
                            tt="")

# From Flask documentation: https://flask.palletsprojects.com/en/1.1.x/patterns/fileuploads/
@app.route('/upload', methods = ['POST'])
def upload():
    if request.method != 'POST':
        return
    # check if the post request has the file part
    if 'file' not in request.files:
        return redirect('/')
    file = request.files['file']
    # if user does not select file, browser also
    # submit an empty part without filename
    if file.filename == '':
        return redirect('/')
    print(file.filename[-3:])
    if file and file.filename[-3:] == "png":
        filename = secure_filename("original.png")
        file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))

        with open(filename, "rb") as img_file:
            b64Image = base64.b64encode(img_file.read())
        b64Image = "data:image/png;base64," + str(b64Image)[2:-3]
        im = Image.open(filename)
        w, h = im.size
        return render_template('edited.html', 
                                data=b64Image,
                                r=0, 
                                g=0, 
                                b=0,
                                c=0,
                                w=w,
                                h=h,
                                th=1,
                                hue=0,
                                sat=0,
                                light=0,
                                trgb="",
                                thsl="",
                                tt="")