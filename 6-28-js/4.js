class Animal {
    constructor(name){
        this.name = name;
    }
    speak(){
        console.log(`${this.name}的叫声：`);
    }
}

class Dog extends Animal {
    constructor(name){
        super(name);
    }
    bark(){
        console.log("Woof!");
    }
}

const mydog = new Dog("blues");
mydog.speak();
mydog.bark();