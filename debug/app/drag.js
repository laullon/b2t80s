function drag() {
    var separator;
    var left;
    var right;
    var main;
    var draging;


    function onMouseDown(e) {
        md = {
            e,
            offsetLeft: separator.offsetLeft,
            offsetTop: separator.offsetTop,
            leftWidth: left.offsetWidth,
            rightWidth: right.offsetWidth
        };

        document.onmousemove = onMouseMove;
        document.onmouseup = () => {
            document.onmousemove = document.onmouseup = null;
        }
    }

    function onMouseUp(e) {
        separator.style.backgroundColor = 'black'
        separator.style.width = "1px";
        draging = false;
    }

    function onMouseMove(e) {
        separator.style.backgroundColor = 'blue'
        separator.style.width = "5px";
        draging = true;
        var delta = {
            x: e.clientX - md.e.clientX,
            y: e.clientY - md.e.clientY
        };
        delta.x = Math.min(Math.max(delta.x, -md.leftWidth), md.rightWidth);
        wl = md.offsetLeft + delta.x
        wr = main.offsetWidth - (1 + wl)

        if ((wl > 160) && (wr > 320)) {
            separator.style.left = wl + "px";
            left.style.width = wl + "px";

            right.style.width = wr + "px";
            right.style.left = (1 + left.offsetWidth) + "px";
        }
        document.getElementById('status').innerHTML = wl + " - " + wr;
    }

    function resize() {
        main.style.width = window.innerWidth + "px"
        right.style.width = (window.innerWidth - (left.offsetWidth + 1)) + "px";
        document.getElementById('status').innerHTML = window.innerWidth + " - " + left.offsetWidth;
    }

    function init() {
        separator = document.getElementById('separator')
        left = document.getElementById('left')
        right = document.getElementById('right')
        main = document.getElementById('main')

        separator.onmousedown = onMouseDown;
        separator.onmouseup = onMouseUp;
        separator.onmouseenter = function () {
            separator.style.backgroundColor = 'blue'
            separator.style.width = "5px";

        }
        separator.onmouseout = function () {
            if (!draging) {
                separator.style.backgroundColor = 'black'
                separator.style.width = "1px";
            }
        }

        document.getElementById('status').innerHTML = "onload done";
        main.style.width = window.innerWidth + "px"
        left.style.width = "160px";
        separator.style.left = "160px";
        right.style.left = "161px";
        right.style.width = (main.offsetWidth - 161) + "px";

        window.onresize = resize
    }

    init()
};