var Client = {};
Client.socket = io('ws://localhost:8032/', {transports: ['websocket']});

Client.sendTest = function(){
    console.log("test sent");
    Client.socket.emit('test');
};

Client.login = function() {
  console.log("login..");
  Client.socket.emit('login', { login: "vasya", password: "password"});
}

Client.askNewPlayer = function(){
    Client.socket.emit('newplayer');
};

Client.sendClick = function(x,y){
  Client.socket.emit('click',{x:x,y:y});
};

Client.socket.on('newplayer',function(data){
    Game.addNewPlayer(data.id, data.point.x, data.point.y);
});

Client.socket.on('allplayers',function(data){
    for(var i = 0; i < data.length; i++){
        Game.addNewPlayer(data[i].id,data[i].x,data[i].y);
    }

    Client.socket.on('move',function(data){
        Game.movePlayer(data.id,data.point.x,data.point.y);
    });

    Client.socket.on('remove',function(id){
        Game.removePlayer(id);
    });
});


