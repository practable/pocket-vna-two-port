/*
 * Remote Lab: RF Switch controller (2-port)
 * Timothy D. Drysdale
 * timothy.d.drysdale@gmail.com
 *
 * Created by Timothy D. Drysdale 25 Jan 2022 - Initial implemention based on github.com/practable/penduino
 * Modified to support 2-Port operation 10 November 2022
 *
 */

 /*********** INSTALLATION ************
  *  Requires https://github.com/khoih-prog/TimerInterrupt
  *  Install using the library manager, or manually, following the readme at the link
  */
  
 /*********** USAGE  *****************
  *  
  *  On startup, monitor at 57600 baud, you will see  {"report":"port","is":"short"}
  *  To change the port, use the commands {"set":"port","to":<port>} where port is
  *  "short", "open", "load", "dut" e.g.
  *  {"set":"port","to":"short"} 
  *  {"set":"port","to":"open"} 
  *  {"set":"port","to":"load"} 
  *  {"set":"port","to":"thru"} 
  *  {"set":"port","to":"dut1"} 
  *  {"set":"port","to":"dut2"} 
  *  {"set":"port","to":"dut3"} 
  *  {"set":"port","to":"dut4"} 
  *  You will get a confirmation message everytime you change the port setting
  *  {"report":"port","is":"short"}
  *  {"report":"port","is":"open"}
  *  {"report":"port","is":"load"}  
  *  {"report":"port","is":"dut1"}  
  *  etc  
  */


//=============================================================
// SET BOTH THESE TO FALSE BEFORE ROLLING OUT IN PRODUCTION

// report additional information (may affect performance)
bool debug = false;
bool trace = false;

//=============================================================

/********** HEADERS ******************/
#include "ArduinoJson-v6.9.1.h"
// See Timer headers below - these need a #define before them


/********** RF SWITCH ***********/
// TODO check which way round power pins are - but both are always on
// so no issue 
#define SWITCH_P1_A      8
#define SWITCH_P1_B      9
#define SWITCH_P1_C      10
#define SWITCH_P1_POWER  3
#define SWITCH_P2_A      4
#define SWITCH_P2_B      5
#define SWITCH_P2_C      6
#define SWITCH_P2_POWER  2

#define PORT_LED     15
#define POWER_LED    16

/********* RF PORTS ************/
// Check the RF output number on the RF switch 
// and see where each standard/dut is connected
// The value in the enum is zero-indexed 
// The value of the port is one-indexed
// So Port 2 in the datasheet is port 1 in the enum etc
enum port1 {
  PORT1_SHORT = 1,
  PORT1_OPEN =  2,
  PORT1_LOAD =  3,
  PORT1_THRU =  4,
  PORT1_DUT1 =  5,
  PORT1_DUT2 =  6,
  PORT1_DUT3 =  7,
  PORT1_DUT4 =  0,
};

enum port2 {
  PORT2_SHORT = 4,
  PORT2_OPEN =  5,
  PORT2_LOAD =  6,
  PORT2_THRU =  3,
  PORT2_DUT1 =  2,
  PORT2_DUT2 =  1,
  PORT2_DUT3 =  0,
  PORT2_DUT4 =  7,
};


/********* RF PORT NAME ************/
static const char name_short[] = "short";
static const char name_open[] = "open";
static const char name_load[] = "load";
static const char name_thru[] = "thru";
static const char name_dut1[] = "dut1";
static const char name_dut2[] = "dut2";
static const char name_dut3[] = "dut3";
static const char name_dut4[] = "dut4";

/*********** LED DISPLAY ***********/
#define LED_SWITCH 13

//this sets the blink code
// do not edit to suit actual RF ports
// display what port was last requested
enum blinkCode {
  BLINK_SHORT = 1,
  BLINK_OPEN = 2,
  BLINK_LOAD = 3,
  BLINK_THRU = 4,
  BLINK_DUT1 = 5,
  BLINK_DUT2 = 6,
  BLINK_DUT3 = 7,
  BLINK_DUT4 = 8,  
};

/********** TIMER *****************/
// must put #define before the include; see https://github.com/khoih-prog/TimerInterrupt/discussions/21
#define USE_TIMER_1     true

// These two includes to be included only in the .ino with setup() to avoid `Multiple Definitions` Linker Error
#include "TimerInterrupt.h"           //https://github.com/khoih-prog/TimerInterrupt
#include "ISR_Timer.h"                //https://github.com/khoih-prog/TimerInterrupt

#define TIMER_INTERVAL_MS        250L



/*********** JSON SERIAL ***********/
#define COMMAND_SIZE 128 
#define REPORT_SIZE 128 
char command[COMMAND_SIZE];
char writeBuffer[REPORT_SIZE];
StaticJsonDocument<COMMAND_SIZE> doc;


/******** OTHER GLOBAL VARIABLES ********/
bool writing;//for serial semaphore
long int count; //counter for displaying port set as blinks.
int blink; //current port state to blink according to blink enum

// pins struct represents RF switch control pin settings
struct pins {
    bool a, b, c;
};

typedef struct pins Pins;

//=============================================================
// Function Prototypes
//=============================================================

void setRFPort(int port1, int port2);
void reportRFPort(const char *name); //const since not modifying the string
void requestSerial(void);
void releaseSerial(void);
bool blinkState(long int count, int blink);
Pins getPins(int port);
void serialPrintPins(int port1, int port2);
/**
 * Defines the valid states for the state machine
 *
 */
typedef enum
{
  STATE_SHORT_BEFORE,
  STATE_SHORT_DURING,
  STATE_OPEN_BEFORE,
  STATE_OPEN_DURING,
  STATE_LOAD_BEFORE,
  STATE_LOAD_DURING,
  STATE_THRU_BEFORE,
  STATE_THRU_DURING,  
  STATE_DUT1_BEFORE,
  STATE_DUT1_DURING,
  STATE_DUT2_BEFORE,
  STATE_DUT2_DURING,
  STATE_DUT3_BEFORE,
  STATE_DUT3_DURING,
  STATE_DUT4_BEFORE,
  STATE_DUT4_DURING,  
 } StateType;

//state Machine function prototypes
//these are the functions that run whilst in each respective state.
void stateShortBefore(void);
void stateShortDuring(void);
void stateOpenBefore(void);
void stateOpenDuring(void);
void stateLoadBefore(void);
void stateLoadDuring(void);
void stateThruBefore(void);
void stateThruDuring(void);
void stateDUT1Before(void);
void stateDUT1During(void);
void stateDUT2Before(void);
void stateDUT2During(void);
void stateDUT3Before(void);
void stateDUT3During(void);
void stateDUT4Before(void);
void stateDUT4During(void);

/**
 * Type definition used to define the state
 */

typedef struct
{
  StateType State; /**< Defines the command */
  void (*func)(void); /**< Defines the function to run */
} StateMachineType;

/**
 * A table that defines the valid states of the state machine and
 * the function that should be executed for each state
 */
StateMachineType StateMachine[] =
{
  {STATE_SHORT_BEFORE, stateShortBefore},
  {STATE_SHORT_DURING, stateShortDuring},
  {STATE_OPEN_BEFORE, stateOpenBefore},
  {STATE_OPEN_DURING, stateOpenDuring}, 
  {STATE_LOAD_BEFORE, stateLoadBefore},
  {STATE_LOAD_DURING, stateLoadDuring},
  {STATE_THRU_BEFORE, stateThruBefore},
  {STATE_THRU_DURING, stateThruDuring},  
  {STATE_DUT1_BEFORE, stateDUT1Before},  
  {STATE_DUT1_DURING, stateDUT1During},  
  {STATE_DUT2_BEFORE, stateDUT2Before},  
  {STATE_DUT2_DURING, stateDUT2During},  
  {STATE_DUT3_BEFORE, stateDUT3Before},  
  {STATE_DUT3_DURING, stateDUT3During},  
  {STATE_DUT4_BEFORE, stateDUT4Before},  
  {STATE_DUT4_DURING, stateDUT4During},        
};

int numStates = 16;

/**
 * Stores the current state of the state machine
 */

StateType state = STATE_SHORT_BEFORE;    //Start with the first cal standard

//============================================================================
//           DEFINE STATE MACHINE FUNCTIONS
//
// A typical sequence uses the BEFORE state as the entry point for a given task:
//
//    STATE_<STATE>_BEFORE
//    STATE_<STATE>_DURING (repeat until finished)
//    STATE_<STATE>_AFTER
//
//  Where not currently required, these states are provided as placeholders
//  for ease of future modification.
//
//  If a change to parameters or command is needed, the typical sequence becomes 
//
//    STATE_<STATE>_BEFORE
//    STATE_<STATE>_DURING (repeat until change)
//    STATE_<STATE>_CHANGE_<SOMETHING>
//    STATE_<STATE>_DURING (repeat until finished)
//    STATE_<STATE>_AFTER
//
//  Note that you should NOT go back to BEFORE following a change -
//  reimplement any required logic from BEFORE in CHANGE_<SOMETHING>
//  so that truly one-shot setup stuff can be put in BEFORE.
// 
//  The default NEXT state is set in the first line of each state,
//  then overridden if need be.
//
//  An exception to this convention are the additional STATE_POSITION_WAITING,
//  and STATE_POSITION_READY, which are intended to let you set the PID parameters
//  before running the PID control routine, so that users can recover from setting
//  large integral coefficients, which cause an automatic stop for going out
//  of bounds before a lower PID value can be set. On entering position PID mode
//  the before state zeroes, the disk, then passes to the READY state.
//
//============================================================================

void stateShortBefore(void) {

  state = STATE_SHORT_DURING;
  setRFPort(PORT1_SHORT, PORT2_SHORT);
  reportRFPort(name_short);
  blink = BLINK_SHORT;
}

void stateOpenBefore(void) {

  state = STATE_OPEN_DURING;
  setRFPort(PORT1_OPEN, PORT2_OPEN);
  reportRFPort(name_open);
  blink = BLINK_OPEN;
}

void stateLoadBefore(void) {

  state = STATE_LOAD_DURING;
  setRFPort(PORT1_LOAD,PORT2_LOAD);
  reportRFPort(name_load);
  blink = BLINK_LOAD;
}


void stateThruBefore(void) {

  state = STATE_THRU_DURING;
  setRFPort(PORT1_THRU,PORT2_THRU);
  reportRFPort(name_thru);
  blink = BLINK_THRU;
}


void stateDUT1Before(void) {

  state = STATE_DUT1_DURING;
  setRFPort(PORT1_DUT1, PORT2_DUT1);
  reportRFPort(name_dut1);
  blink = BLINK_DUT1;
}

void stateDUT2Before(void) {

  state = STATE_DUT2_DURING;
  setRFPort(PORT1_DUT2, PORT2_DUT2);
  reportRFPort(name_dut2);
  blink = BLINK_DUT2;
}

void stateDUT3Before(void) {

  state = STATE_DUT3_DURING;
  setRFPort(PORT1_DUT3, PORT2_DUT3);
  reportRFPort(name_dut3);
  blink = BLINK_DUT3;
}

void stateDUT4Before(void) {

  state = STATE_DUT4_DURING;
  setRFPort(PORT1_DUT4, PORT2_DUT4);
  reportRFPort(name_dut4);
  blink = BLINK_DUT4;
}
void stateShortDuring(void) {

  state = STATE_SHORT_DURING;
  // do nothing

}

void stateOpenDuring(void) {

  state = STATE_OPEN_DURING;
  // do nothing
}

void stateLoadDuring(void) {

  state = STATE_LOAD_DURING;
  // do nothing
}

void stateThruDuring(void) {

  state = STATE_THRU_DURING;
  // do nothing
}

void stateDUT1During(void) {

  state = STATE_DUT1_DURING;
  // do nothing
}

void stateDUT2During(void) {

  state = STATE_DUT2_DURING;
  // do nothing
}

void stateDUT3During(void) {

  state = STATE_DUT3_DURING;
  // do nothing
}

void stateDUT4During(void) {

  state = STATE_DUT4_DURING;
  // do nothing
}



void setRFPort(int port1, int port2){

  Pins p1, p2;
  p1 = getPins(port1);
  p2 = getPins(port2);

  digitalWrite(SWITCH_P1_A, p1.a);
  digitalWrite(SWITCH_P1_B, p1.b);
  digitalWrite(SWITCH_P1_C, p1.c);
  digitalWrite(SWITCH_P2_A, p2.a);
  digitalWrite(SWITCH_P2_B, p2.b);
  digitalWrite(SWITCH_P2_C, p2.c);
  
}

void reportRFPort(const char *name){
  requestSerial();
  Serial.print("{\"report\":\"port\",\"is\":\"");
  Serial.print(name);
  Serial.println("\"}");
  releaseSerial();
  
}

bool blinkState(long int count, int blink){

  // we want a pattern like this
  // SHORT: (10)(00)(00)(00)(00)(00)(00)(00)(delay)
  // OPEN:  (10)(10)(00)(00)(00)(00)(00)(00)(delay)
  // LOAD:  (10)(10)(10)(00)(00)(00)(00)(00)(delay)
  // THRU:  (10)(10)(10)(10)(00)(00)(00)(00)(delay)
  // DUT1:  (10)(10)(10)(10)(10)(00)(00)(00)(delay) 
  // DUT2:  (10)(10)(10)(10)(10)(10)(00)(00)(delay) 
  // DUT3:  (10)(10)(10)(10)(10)(10)(10)(00)(delay) 
  // DUT4:  (10)(10)(10)(10)(10)(10)(10)(10)(delay)   
  // where delay is the same length again.

  switch (count%24){
      case 0:
        return true;
      case 2:
        switch (blink){
            case BLINK_OPEN:
            case BLINK_LOAD:
            case BLINK_THRU:
            case BLINK_DUT1:
            case BLINK_DUT2:
            case BLINK_DUT3:
            case BLINK_DUT4:            
              return true;            
            default:
              return false;  
        }
      case 4:
        switch (blink){
            case BLINK_LOAD:
            case BLINK_THRU:
            case BLINK_DUT1:
            case BLINK_DUT2:
            case BLINK_DUT3:
            case BLINK_DUT4: 
              return true;            
            default:
              return false;  
        }  
       case 6:
        switch (blink){
            case BLINK_THRU:
            case BLINK_DUT1:
            case BLINK_DUT2:
            case BLINK_DUT3:
            case BLINK_DUT4: 
              return true;            
            default:
              return false;  
        } 
       case 8:
        switch (blink){
            case BLINK_DUT1:
            case BLINK_DUT2:
            case BLINK_DUT3:
            case BLINK_DUT4: 
              return true;            
            default:
              return false;  
        } 
       case 10:
        switch (blink){
            case BLINK_DUT2:
            case BLINK_DUT3:
            case BLINK_DUT4: 
              return true;            
            default:
              return false;  
        }      
       case 12:
        switch (blink){
            case BLINK_DUT3:
            case BLINK_DUT4: 
              return true;            
            default:
              return false;  
        }  
       case 14:
        switch (blink){
            case BLINK_DUT4: 
              return true;            
            default:
              return false;  
        }                             
      default: 
        return false;
  }
  

  
}

Pins getPins(int port) {
 
  Pins p;

  p.c = (port & ( 1 << 0 )) >> 0;
  p.b = (port & ( 1 << 1 )) >> 1;
  p.a = (port & ( 1 << 2 )) >> 2;

  return p;
    
}

  
//*****************************************************************************

//STATE MACHINE RUN FUNCTION
void SMRun(void)
{
  if (state < numStates)
  {
    state = readSerialJSON(state);      // check for incoming commands received
    (*StateMachine[state].func)();        //reads the current state and then runs the associated function

  }
  else {
    Serial.println("Exception in State Machine");
  }

}


//This function is run on a timer interrupt 
void TimerHandler()
{
    digitalWrite(PORT_LED,blinkState(count, blink));
    count++;
    
    digitalWrite(POWER_LED, digitalRead(SWITCH_P1_POWER)&&digitalRead(SWITCH_P2_POWER));
}

//===================================================================================
//====================== SETUP AND LOOP =============================================
//===================================================================================

void setup() {

  // pins for RF swith power
  pinMode(SWITCH_P1_POWER, OUTPUT);  
  pinMode(SWITCH_P2_POWER, OUTPUT);     
    
  // pins for the RF switch control
  pinMode(SWITCH_P1_A, OUTPUT);
  pinMode(SWITCH_P1_B, OUTPUT);
  pinMode(SWITCH_P1_C, OUTPUT);
  pinMode(SWITCH_P2_A, OUTPUT);
  pinMode(SWITCH_P2_B, OUTPUT);
  pinMode(SWITCH_P2_C, OUTPUT);
    
  // pins for LED display
  pinMode(PORT_LED, OUTPUT);
  pinMode(POWER_LED, OUTPUT);

  // Turn on the RF switches
  digitalWrite(SWITCH_P1_POWER,1);
  digitalWrite(SWITCH_P2_POWER,1);

  // Init timer ITimer1
  ITimer1.init();

  // timer
  ITimer1.attachInterruptInterval(TIMER_INTERVAL_MS, TimerHandler);
  
  Serial.setTimeout(50);
  Serial.begin(57600);
 
  while (! Serial); //wait for serial to start 

  count = 0; //initialise counter for use in blink code

  if (debug){
    /*
    #define SWITCH_P1_A      8
#define SWITCH_P1_B      9
#define SWITCH_P1_C      10
#define SWITCH_P1_POWER  3
#define SWITCH_P2_A      4
#define SWITCH_P2_B      5
#define SWITCH_P2_C      6
*/
  Serial.print("       4 5 6   8 9 10\n");
  Serial.print("        P2      P1     P2  P1\n");
  Serial.print("short: ");
  serialPrintPins(PORT1_SHORT, PORT2_SHORT);
  Serial.print("\nopen : ");
  serialPrintPins(PORT1_OPEN, PORT2_OPEN);  
  Serial.print("\nload : ");
  serialPrintPins(PORT1_LOAD, PORT2_LOAD);  
  Serial.print("\nthru : ");
  serialPrintPins(PORT1_THRU, PORT2_THRU);  
  Serial.print("\ndut1 : ");
  serialPrintPins(PORT1_DUT1, PORT2_DUT1);      
  Serial.print("\ndut2 : ");
  serialPrintPins(PORT1_DUT2, PORT2_DUT2);  
  Serial.print("\ndut3 : ");
  serialPrintPins(PORT1_DUT3, PORT2_DUT3);  
  Serial.print("\ndut4 : ");
  serialPrintPins(PORT1_DUT4, PORT2_DUT4);  
  Serial.print("\n");
  }
}

void serialPrintPins(int port1, int port2){
  Pins p1, p2;

  p1 = getPins(port1);
  p2 = getPins(port2);
  Serial.print(p2.a);
  Serial.print(" ");
  Serial.print(p2.b);
  Serial.print(" ");
  Serial.print(p2.c);
  Serial.print("   ");  
  Serial.print(p1.a);
  Serial.print(" ");
  Serial.print(p1.b);
  Serial.print(" ");
  Serial.print(p1.c);
  Serial.print("   ");
  Serial.print(port1);
  Serial.print("   ");
  Serial.print(port2);
  
}

void loop() {

  // update state machine (which will also run tasks flagged by interrupts)
  SMRun();

}

//===================================================================================
//====================== SUPPORTING FUNCTIONS========================================
//
//
//
//            These may be used by one or more states
//
//
//
//
//===================================================================================

float hertzFromSeconds(float t) {
  return 1.0f / t;
}


float hertzFromMillis(float t) {
  return 1000.0f / t;
}

float secondsFromMillis(float t) {
  return t / 1000.0f;
}

float secondsFromMicros(float t) {
  return t / 1000000.0f;
}


//===================================================================================
//======================  READ AND PARSE JSON COMMMANDS =============================
//
//  This function can and does change the state of the state machine
//
//===================================================================================

StateType readSerialJSON(StateType state) {

  if(Serial.available() > 0) {

    Serial.readBytesUntil(10, command, COMMAND_SIZE);
    deserializeJson(doc, command);

    const char* set = doc["set"];

    if(strcmp(set, "port")==0) {
 
        const char* port = doc["to"];

        if(strcmp(port, name_short) == 0) {
          state = STATE_SHORT_BEFORE;
        }
        else if(strcmp(port, name_open) == 0) {
          state = STATE_OPEN_BEFORE;
        }
        else if(strcmp(port, name_load) == 0) {
          state = STATE_LOAD_BEFORE;
        }    
        else if(strcmp(port, name_thru) == 0) {
          state = STATE_THRU_BEFORE;
        }             
        else if(strcmp(port, name_dut1) == 0) {
          state = STATE_DUT1_BEFORE;
        }          
        else if(strcmp(port, name_dut2) == 0) {
          state = STATE_DUT2_BEFORE;
        } 
        else if(strcmp(port, name_dut3) == 0) {
          state = STATE_DUT3_BEFORE;
        } 
        else if(strcmp(port, name_dut4) == 0) {
          state = STATE_DUT4_BEFORE;
        } 
    }
  }
 
  return state;     //return whatever state it changed to or maintain the state.
}

void requestSerial(void){
  while(writing); //wait for port to free up
  writing = true;
}

void releaseSerial(void){
  writing = false;
}


 
