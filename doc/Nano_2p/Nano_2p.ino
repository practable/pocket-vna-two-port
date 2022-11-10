  
 
#define LED 13          // Pin 13 is connected to the LED
char rxChar= 0;         // RXcHAR holds the received command.

//=== function to print the command list:  ===========================
void printHelp(void){
  Serial.println("System on");
  }
  
//---------------- setup ---------------------------------------------
void setup(){
  Serial.begin(9600);   // Open serial port (9600 bauds).
  pinMode(LED, OUTPUT); // Sets pin 13 as OUTPUT.

  
   pinMode(2, OUTPUT);    // pin D2 ->   to 5 V
  pinMode(3, OUTPUT);    // pin D3  -> to 5 V


pinMode(4, OUTPUT);   // Ouput A coupler 1
pinMode(5, OUTPUT);   // Ouput B coupler 1 
pinMode(3, OUTPUT);   // Ouput C coupler 1


pinMode(8, OUTPUT);   // Ouput A coupler 2
pinMode(9, OUTPUT);   // Ouput B coupler 2 
pinMode(10, OUTPUT);   // Ouput C coupler 2


pinMode(16, OUTPUT);    // LED for coupler 1 active
pinMode(15, OUTPUT);    // LED for coupler 2 active

  
  Serial.flush();       // Clear receive buffer.
  printHelp();          // Print the command list.
}

//--------------- loop ----------------------------------------------- 
void loop(){
    
      digitalWrite(2, HIGH);;    //Power to coupler 1
      digitalWrite(3, HIGH);    // power to Coupler 2

 //digitalWrite(4, HIGH);    // Ouput A coupler 1
// digitalWrite(5, HIGH);    // Ouput B coupler 1
//digitalWrite(6, HIGH);    // Ouput C coupler 1




 //digitalWrite(8, HIGH);    // Ouput A coupler 1
// digitalWrite(9, HIGH);    // Ouput B coupler 1
// digitalWrite(10, HIGH);    // Ouput C coupler 1
 
 //  digitalWrite(16, HIGH);    // LED for coupler 1 active
  //digitalWrite(15, HIGH);    // LED for coupler 2 active

   
  if (Serial.available() >0){          // Check receive buffer.
    rxChar = Serial.read();            // Save character received. 
    Serial.flush();                    // Clear receive buffer.


  switch (rxChar) {





    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 4 5 6
    case '1':                         
 {       

         digitalWrite(4, HIGH);
          digitalWrite(5, LOW);
           digitalWrite(6, LOW);

          digitalWrite(16, HIGH);    // LED for coupler 1 active
          digitalWrite(15, LOW);
          
          Serial.println("SHORT SW2");
  }
 
        break;

    case '2':                         
 {       

         digitalWrite(4, HIGH);
          digitalWrite(5, LOW);
           digitalWrite(6, HIGH);

          digitalWrite(16, HIGH);    // LED for coupler 1 active
        digitalWrite(15, LOW);
          
          Serial.println("OPEN SW2");
  }
 
        break;


            case '3':                         
 {       

         digitalWrite(4, HIGH);
          digitalWrite(5, HIGH);
           digitalWrite(6, LOW);

          digitalWrite(16, HIGH);    // LED for coupler 1 active
          digitalWrite(15, LOW);
          
          Serial.println("LOAD SW2");
  }
 
        break;



////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
        

                    case '4':                         
 {       

         digitalWrite(8, LOW);
          digitalWrite(9, LOW);
           digitalWrite(10, HIGH);

          digitalWrite(15, HIGH);    // LED for coupler 2 active
          digitalWrite(16, LOW);
          
          Serial.println("SHORT SW1");
  }
 
        break;


                

                    case '5':                         
 {       

         digitalWrite(8, LOW);
          digitalWrite(9, HIGH);
           digitalWrite(10, LOW);

          digitalWrite(15, HIGH);    // LED for coupler 2 active
         digitalWrite(16, LOW);         
          
          Serial.println("OPEN SW1");
  }
 
        break;


                            case '6':                         
 {       

         digitalWrite(8, LOW);
          digitalWrite(9, HIGH);
           digitalWrite(10, HIGH);

          digitalWrite(15, HIGH);    // LED for coupler 2 active
        digitalWrite(16, LOW);
          
          Serial.println("LOAD SW1");
  }
 
        break;

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


                            case '7':                         
 {       

         digitalWrite(4, LOW);
          digitalWrite(5, HIGH);
           digitalWrite(6, HIGH);


           digitalWrite(8, HIGH);
          digitalWrite(9, LOW);
           digitalWrite(10, LOW);




          digitalWrite(15, HIGH);    // LED for coupler 2 active
          digitalWrite(16, HIGH);    // LED for coupler 1 active


          



          
          Serial.println("THRU");
  }
 
        break;


                            case '8':                         
 {       

         digitalWrite(4, LOW);
          digitalWrite(5, HIGH);
           digitalWrite(6, LOW);


           digitalWrite(8, HIGH);
          digitalWrite(9, LOW);
           digitalWrite(10, HIGH);




          digitalWrite(15, HIGH);    // LED for coupler 2 active
          digitalWrite(16, HIGH);    // LED for coupler 1 active


          



          
          Serial.println("DUT 1");
  }
 
        break;

            ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////THOSE ARE THE ONE PORT CASES FOR SWITCH 2, CAN BE CHANGED!

                                case '9':                         
 {       

         digitalWrite(4, LOW);
          digitalWrite(5, LOW);
           digitalWrite(6, LOW);






          digitalWrite(16, HIGH);    // LED for coupler 2 active
        digitalWrite(15, LOW);
          
          



          
          Serial.println("Port 0 of the SW2 is active (reflection)");
  }
 
        break;


                                        case 'a':                         
 {       

         digitalWrite(4, LOW);
          digitalWrite(5, LOW);
           digitalWrite(6, HIGH);





 

          digitalWrite(16, HIGH);    // LED for coupler 2 active
        digitalWrite(15, LOW);
          
          



          
          Serial.println("Port 1 of the SW2 is active (reflection)");
  }
 
        break;


        
                                        case 'b':                         
 {       

         digitalWrite(4, HIGH);
          digitalWrite(5, HIGH);
           digitalWrite(6, HIGH);





 

          digitalWrite(16, HIGH);    // LED for coupler 2 active
        digitalWrite(15, LOW);
          
          



          
          Serial.println("Port 7 of the SW2 is active (reflection)");
  }
 
        break;

        
            ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////THOSE ARE THE ONE PORT CASES FOR SWITCH 1, CAN BE CHANGED!
           
                                        case 'c':                         
 {       

       digitalWrite(8, LOW);
          digitalWrite(9, LOW);
           digitalWrite(10, LOW);

          digitalWrite(15, HIGH);    // LED for coupler 1 active
          digitalWrite(16, LOW);
          
          Serial.println("port 0 for switch 1 is active (reflection)");

  }
 
        break;


                                                case 'd':                         
 {       

       digitalWrite(8, HIGH);
          digitalWrite(9, HIGH);
           digitalWrite(10,LOW);

          digitalWrite(15, HIGH);    // LED for coupler 1 active
          digitalWrite(16, LOW);
          
          Serial.println("port 6 for switch 1 is active (reflection)");

  }
 
        break;

                                                        case 'e':                         
 {       

       digitalWrite(8, HIGH);
          digitalWrite(9, HIGH);
           digitalWrite(10,HIGH);

          digitalWrite(15, HIGH);    // LED for coupler 1 active
          digitalWrite(16, LOW);
          
          Serial.println("port 7 for switch 1 is active (reflection)");

  }
 
        break;

                    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////Special case for DUT 2 if lab will be for this

                                                case 'f':                         
 {       

       digitalWrite(8, HIGH);
          digitalWrite(9, HIGH);
           digitalWrite(10,LOW);

            digitalWrite(4, LOW);
          digitalWrite(5, LOW);
           digitalWrite(6, HIGH);

           

          digitalWrite(15, HIGH);    // LED for coupler 1 active
          digitalWrite(16, HIGH);  // LED for coupler 2 active
          
          Serial.println("port 6 for switch 1 and port 1 for switch 2 is active (DUT2)");

  }
 
        break;

                                                case 'g':                         
 {       

       digitalWrite(8, HIGH);
          digitalWrite(9, HIGH);
           digitalWrite(10,HIGH);

            digitalWrite(4, LOW);
          digitalWrite(5, LOW);
           digitalWrite(6, LOW);

           

          digitalWrite(15, HIGH);    // LED for coupler 1 active
          digitalWrite(16, HIGH);  // LED for coupler 2 active
          
          Serial.println("port 7 for switch 1 and port 0 for switch 2 is active (DUT3)");

  }
 
        break;
                                                        case 'h':                         
 {       

       digitalWrite(8, LOW);
          digitalWrite(9, LOW);
           digitalWrite(10,LOW);

            digitalWrite(4, HIGH);
          digitalWrite(5, HIGH);
           digitalWrite(6, HIGH);

           

          digitalWrite(15, HIGH);    // LED for coupler 1 active
          digitalWrite(16, HIGH);  // LED for coupler 2 active
          
          Serial.println("port 0 for switch 1 and port 7 for switch 2 is active (DUT4)");

  }
 
        break;
    }
  }
}
// End of the Sketch.
