const fs = require('fs');

const express = require('express')
    , http = require('http')
    , app = express()
    , server = http.createServer(app);
const path = require("path");


app.use(express.static(path.join(__dirname, './public')));

app.get('/', function (req, res) {
    res.writeHead(200, {"Content-Type": "text/html"});
    fs.createReadStream("./index.html").pipe(res);
})

app.get('/register', function (req, res) {
    res.writeHead(200, {"Content-Type": "text/html"});
    fs.createReadStream("./register.html").pipe(res);
})

app.get('/dashboard/total', function (req, res) {
    res.writeHead(200, {"Content-Type": "text/html"});
    fs.createReadStream("./dashboard/main.html").pipe(res);
})

app.get('/dashboard/detail', function (req, res) {
    res.writeHead(200, {"Content-Type": "text/html"});
    fs.createReadStream("./dashboard/project-detail.html").pipe(res);
})

app.get('/project/progress', function (req, res) {
    res.writeHead(200, {"Content-Type": "text/html"});
    fs.createReadStream("./project/project.html").pipe(res);
})
//
// app.get('/project/complete', function (req, res) {
//     res.writeHead(200, {"Content-Type": "text/html"});
//     fs.createReadStream("./project/project-ed.html").pipe(res);
// })
// app.get('/project/ing', function (req, res) {
//     res.writeHead(200, {"Content-Type": "text/html"});
//     fs.createReadStream("./project/project-ing.html").pipe(res);
// })

app.get('/project/create', function (req, res) {
    res.writeHead(200, {"Content-Type": "text/html"});
    fs.createReadStream("./project/project-making.html").pipe(res);
})

app.get('/setting/user', function (req, res) {
    res.writeHead(200, {"Content-Type": "text/html"});
    fs.createReadStream("./setting/person-setting.html").pipe(res);
})


app.get('/setting/smtp', function (req, res) {
    res.writeHead(200, {"Content-Type": "text/html"});
    fs.createReadStream("./setting/smtp-setting.html").pipe(res);
})


app.get('/manager/target', function (req, res) {
    res.writeHead(200, {"Content-Type": "text/html"});
    fs.createReadStream("./target/target.html").pipe(res);
})

app.get('/manager/template', function (req, res) {
    res.writeHead(200, {"Content-Type": "text/html"});
    fs.createReadStream("./templates/template.html").pipe(res);
})




app.get('/setting')


// 맨아래 둬야하는건가? 쉣인데?
app.use(function (req, res, next) {
    res.status(404);
    fs.createReadStream("./error.html").pipe(res);
});


server.listen(8080, function () {
    console.log('Express server listening on port ' + server.address().port);
})
