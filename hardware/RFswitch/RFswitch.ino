
 
#define LED 13          // Pin 13 is connected to the LED
char rxChar= 0;         // RXcHAR holds the received command.

//=== function to print the command list:  ===========================
void printHelp(void){
  Serial.println("1 for RF1, 2 for RF2 , 3 for RF3, 4 for RF4");
  }
  
//---------------- setup ---------------------------------------------
void setup(){
  Serial.begin(9600);   // Open serial port (9600 bauds).
  pinMode(LED, OUTPUT); // Sets pin 13 as OUTPUT.

   pinMode(8, OUTPUT);    // pin 8
   pinMode(12, OUTPUT);    // sets the digital pin 12 as output

  
  Serial.flush();       // Clear receive buffer.
  printHelp();          // Print the command list.
}

//--------------- loop ----------------------------------------------- 
void loop(){
  if (Serial.available() >0){          // Check receive buffer.
    rxChar = Serial.read();            // Save character received. 
    Serial.flush();                    // Clear receive buffer.
  
  switch (rxChar) {





    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    case '1':                         
 {        // If LED is Off:
          digitalWrite(LED,HIGH);     // Turn On the LED.

          digitalWrite(8, LOW);
          digitalWrite(12, LOW);

       
          Serial.println("RF1 ON");
  }
 
        break;


    case '2':                         
 {        // If LED is Off:
          digitalWrite(LED,LOW);     // Turn Off the LED.

          digitalWrite(8, LOW);
          digitalWrite(12, HIGH);

       
          Serial.println("RF2 ON");
  }
 
        break;

            case '3':                         
 {        // If LED is Off:
          digitalWrite(LED,LOW);     // Turn Off the LED.

          digitalWrite(8, HIGH);
          digitalWrite(12, LOW);

       
          Serial.println("RF3 ON");
  }
 
        break;


            case '4':                          
 {        // If LED is Off:
          digitalWrite(LED,HIGH);     // Turn On the LED.

          digitalWrite(8, HIGH);
          digitalWrite(12, HIGH);

       
          Serial.println("RF4 ON");
  }
 
        break;



////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
        

    }
  }
}
// End of the Sketch.
