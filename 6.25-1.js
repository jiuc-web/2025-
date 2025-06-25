function checkIfInstanceOf(obj, classFunction) {
  if(obj === null || classFunction === null || obj === undefined || classFunction === undefined) {
     return false;
  }
   while(obj.__proto__ !== null) {
         if(obj.__proto__ === classFunction.prototype) {
            return true;
         } else {
            obj = obj.__proto__;
         }
   }

  return false;
};