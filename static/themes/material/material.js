$('.top-icon').each(function(index, element) {
    $(element).css('background-color','hsl('+getRandomInt(0, 360)+',30%,70%)');
});

function getRandomInt (min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

//TODO: Use jquery to move left-hand navbar icons into their own dropdown