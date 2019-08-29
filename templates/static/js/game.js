document.onreadystatechange = function () {
    var state = -1;
    var isLeft = true;

    var playerPos = 20;
    var contPos = 773;
    var textPos = 20;
    
    var left = 0;
    var right = 0;
    var leftW = 0;
    var rightW = 0;

    if (document.readyState == "complete") {

        var game = new PixelJS.Engine();
        game.init({
            container: 'game_container',
            width: 800,
            height: 735
        });
        
        //layers
        var backgroundLayer = game.createLayer('background');
        var playerLayer = game.createLayer('players')
        var ballLayer = game.createLayer('ball');
        var contLayer = game.createLayer("contestant");
        var wallsLayer = game.createLayer("walls");
        var textLayer = game.createLayer("text");
        var playerPosLayer = game.createLayer("playerPosText");

        //layers propeties
        textLayer.redraw = true;
        playerPosLayer.redraw = false;
        wallsLayer.zIndex = -3;
        backgroundLayer.static = true;
        ballLayer.static = false;

        //sounds
        var wallHitSound = game.createSound('wallHitSound');
        wallHitSound.prepare({ name: 'newWallHit.mp3' });
        var playerHitSound = game.createSound('playerHitSound');
        playerHitSound.prepare({ name: 'playerHit.mp3' });

        
        //main code
        var field = backgroundLayer.createEntity();
        field.pos = { x: 0, y: 0 };
        field.asset = new PixelJS.Sprite();
        field.asset.prepare({
            name: 'field.png',
            size: { 
                width: 800, 
                height: 800
            }
        });

        var topWall = wallsLayer.createEntity();
        topWall.isCollidable = true;
        topWall.pos = { x: 5, y: 130};
        topWall.size = {width: 787, height: 5};
        topWall.asset = new PixelJS.Sprite();
        topWall.asset.prepare({
            name: 'wall.png',
            size: { 
                width: 787, 
                height: 5
            }
        });
        wallsLayer.registerCollidable(topWall);

        var bottomWall = wallsLayer.createEntity();
        bottomWall.isCollidable = true;
        bottomWall.pos = { x: 5, y: 725};
        bottomWall.size = {width: 787, height: 5};
        bottomWall.asset = new PixelJS.Sprite();
        bottomWall.asset.prepare({
            name: 'wall.png',
            size: { 
                width: 787, 
                height: 5
            }
        });
        wallsLayer.registerCollidable(bottomWall);

        var player = new PixelJS.Player();
        player.addToLayer(playerLayer);
        player.isCollidable = true;
        player.pos = {x: playerPos, y: 386};
        player.size = {width: 7, height: 90};
        player.velocity = {x: 0, y: 450};
        player.isAnimatedSprite = false;
        player.asset = new PixelJS.Sprite();
        player.asset.prepare({ 
            name: 'brightPlayer.png',
            size: { 
                width: 7, 
                height: 90
            }
        });

        var contestant = contLayer.createEntity();
        contestant.isCollidable = true;
        contestant.pos = {x: contPos, y: 385};
        contestant.size = {width: 7, height: 90};
        contestant.asset = new PixelJS.Sprite();
        contestant.asset.prepare({
            name: 'brightPlayer.png',
            size: { 
                width: 7, 
                height: 90
            }
        });

        playerLayer.registerCollidable(player);
        contLayer.registerCollidable(contestant);

        var ball = ballLayer.createEntity();

        ball.isCollidable = true;
        ball.pos = {x: 396, y: 426};
        ball.velocity = {x: 150, y: 150};
        ball.size = {width: 8, height: 8};
        ball.asset = new PixelJS.Sprite();
        ball.asset.prepare({
            name: 'ball.png',
            size: { 
                width: 8, 
                height: 8
            }
        });
        ballLayer.registerCollidable(ball);

        ball.onCollide(function (entity) {
            switch(entity){
                case topWall:
                    wallHitSound.play();
                break;
                case bottomWall:
                    wallHitSound.play();
                break;
                case player:
                    playerHitSound.play();
                break;
                case contestant:
                    playerHitSound.play();
                break;
            }
        });        

        var socket = new WebSocket("ws://" + host + "/game");

        socket.onclose = function (event) {
            alert("Disconnected from server!");
        }

        socket.onmessage = function(event) {
            var message = JSON.parse(event.data, function(key, value) {
                if(key == 'Command') return new String(value);
                if(key == 'IsLeft' || key == 'Scored') return new Boolean(value);
                if(key == 'Player' || key == 'BallX' || key == 'BallY') return parseFloat(value);
                if(key == 'LeftP' || key == 'RightP') return parseInt(value);
                return value;
            });

            switch(message.Command.valueOf()){
                case "init":
                    isLeft = message.IsLeft.valueOf();
                    state = 1;

                    if(!isLeft){
                        player.pos.x = contPos;
                        contestant.pos.x = playerPos;

                        textPos = 670;
                    }else{
                        player.pos.x = playerPos;
                        contestant.pos.x = contPos;
                    }
                    leftW = rightW = 0;
                    left = right = 0;
                break;
                case "ready":
                    state = 0;
                break;
                case "game":
                    contestant.pos.y = message.Player;
                    ball.pos = {x: message.BallX, y: message.BallY};
                break;
                case "score":
                    left = message.LeftP;
                    right = message.RightP;

                    if(left == 5 || right == 5)
                    {
                        if(left == 5) {leftW++;}
                        if(right == 5) {rightW++;}
                        state = 3;
                    }
                break;
                case "opponentLeft":
                    state = 4;
                break;
            }
        }

        function DrawScore() {
            textLayer.drawText("Score", 
              400, 
              35, 
              '20pt "Trebuchet MS", Helvetica, sans-serif', 
              '#ffffff', 
              'center'
            );

            textLayer.drawText(left + " : " + right, 
              400, 
              62, 
              '20pt "Trebuchet MS", Helvetica, sans-serif', 
              '#ffffff', 
              'center'
            );

            textLayer.drawText("Match points", 
              400, 
              88, 
              '20pt "Trebuchet MS", Helvetica, sans-serif', 
              '#ffffff', 
              'center'
            );
            textLayer.drawText(leftW + " : " + rightW, 
			  400, 
			  116, 
			  '20pt "Trebuchet MS", Helvetica, sans-serif', 
			  '#ffffff', 
			  'center'
			);
        }

        function DrawOnScreen(text) {
            textLayer.drawText(text, 
              402.5, 
              75, 
              '20pt "Trebuchet MS", Helvetica, sans-serif', 
              '#ffffff', 
              'center'
            );   
        }

        game.on('keyDown', function (keyCode) {
            if (keyCode === PixelJS.Keys.Enter) {
                if(state == 3){
                    state = 1;

                    left = right = 0;

                    var out = {
                        Command: "readyState"
                    }

                    socket.send(JSON.stringify(out));
                } else if (state == 4) {

                    var out = {
                        Command: "findNew"
                    }

                    socket.send(JSON.stringify(out));        
                }
                
            }
        });

        var timer = 3;

        game.loadAndRun(function (elapsedTime, dt) {
            player.canMoveDown = (player.pos.y + player.size.height) < 725;
            player.canMoveUp = player.pos.y > 135;

            textLayer.drawText("Your side", 
              textPos, 
              116, 
              '20pt "Trebuchet MS", Helvetica, sans-serif', 
              '#ffffff', 
              'left'
            );

            switch(state){
                case -1:
                    DrawOnScreen("Please, refresh the page");
                break;
                case 0:
                    if(timer > 0)
                    {
                        timer -= dt;
                        DrawOnScreen("Game starts in " + (Math.floor(timer) + 1));
                    }
                    else
                    {
                        timer = 3;
                        var info = {
                            Command: "start"
                        }
                        socket.send(JSON.stringify(info));
                        state = 2;
                    }
                break;
                case 1: //waiting for the contestant
                    DrawOnScreen('Waiting for the opponent ...');
                break;
                case 2: //game
                    DrawScore();

                    var output;

                    var info = {
                        Command: "coords",
                        Player: player.pos.y,
                    }

                    output = JSON.stringify(info);

                    socket.send(output);
                break;
                case 3: //win/lose
                    if(isLeft)
                    {
                        if(left == 5)
                            DrawOnScreen('Congrats, you won!(Press Enter to play again)');
                        else
                            DrawOnScreen('You lost, ha-ha!(Press Enter to play again)');
                    }
                    else
                    {
                        if(right == 5)
                            DrawOnScreen('Congrats, you won!(Press Enter to play again)');
                        else
                            DrawOnScreen('You lost, ha-ha!(Press Enter to play again)');
                    }
                break;
                case 4:
                    DrawOnScreen("Your opponent left, press Enter to find new.")
                break;
            }
        });
    }
}


