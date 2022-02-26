 import _ from 'lodash';
 const {User, UserId} = require('../proto/user.api_pb.js');
 const {UserServiceClient} = require('../proto/user.api_grpc_web_pb.js');

 var client = new UserServiceClient('http://localhost:8080');

 var request;
 request = new UserId();
 request.setId(1);



  function component() {
    var element = document.createElement('div');
    client.getUser(request, {}, (err, response) => {
      console.log(err);
      console.log(response.toObject());
      element.innerHTML = _.join(['Hello', response.toObject().name], ' ');
    });
   // Lodash, currently included via a script, is required for this line to work
   // Lodash, now imported by this script
    element.innerHTML = _.join(['Hello', 'Hongshengjie'], ' ');

    return element;
  }

  document.body.appendChild(component());
