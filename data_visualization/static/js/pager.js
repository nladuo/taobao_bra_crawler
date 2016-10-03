$(document).ready(function(){
    $('.pager li>a').mousedown(function () {
        $(this).css("color","#337ab7");
        $(this).css("background-color","#ffffff");
    });

    $('.pager li>a').on("touchstart",function () {
        $(this).css("color","#337ab7");
        $(this).css("background-color","#ffffff");
    });

    $('.pager li>a').mouseup(function () {
        $('.pager li>a').css("background-color","#b4eeb4");
    });

    $('.pager li>a').on("touchend",function () {
        $('.pager li>a').css("background-color","#b4eeb4");
    });
});
