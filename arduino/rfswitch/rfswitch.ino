/*
 * Remote Lab: RF Switch controller
 * Timothy D. Drysdale
 * timothy.d.drysdale@gmail.com
 *
 * Created by Timothy D. Drysdale 25 Jan 2022 - Initial implemention based on github.com/practable/penduino and demo code from Maksim Kuznetcov
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
  *  {"set":"port","to":"dut"} 
  *  You will get a confirmation message everytime you change the port setting
  *  {"report":"port","is":"short"}
  *  {"report":"port","is":"open"}
  *  {"report":"port","is":"load"}  
  *  {"report":"port","is":"dut"}  
  *    
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
// define such that A=low, B=high gives RF Port 2.
#define SWITCH_V1 8
#define SWITCH_V2 12

/********* RF PORTS ************/
// Check the port number on the RF switch 
// and see where each standard/dut is connected
enum port {
  PORT_SHORT = 1,
  PORT_OPEN = 2,
  PORT_LOAD = 3,
  PORT_DUT = 4,
};




/********* RF PORT NAME ************/
static const char name_short[] = "short";
static const char name_open[] = "open";
static const char name_load[] = "load";
static const char name_dut[] = "dut";


/*********** LED DISPLAY ***********/
#define LED_SWITCH 13

//this sets the blink code
// do not edit to suit actual RF ports
// display what port was last requested
enum blinkCode {
  BLINK_SHORT = 1,
  BLINK_OPEN = 2,
  BLINK_LOAD = 3,
  BLINK_DUT = 4,
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

//=============================================================
// Function Prototypes
//=============================================================

void setRFPort(int port);
void reportRFPort(const char *name); //const since not modifying the string
void requestSerial(void);
void releaseSerial(void);
bool blinkState(long int count, int blink);


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
  STATE_DUT_BEFORE,
  STATE_DUT_DURING,
 } StateType;

//state Machine function prototypes
//these are the functions that run whilst in each respective state.
void stateShortBefore(void);
void stateShortDuring(void);
void stateOpenBefore(void);
void stateOpenDuring(void);
void stateLoadBefore(void);
void stateLoadDuring(void);
void stateDUTBefore(void);
void stateDUTDuring(void);

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
  {STATE_DUT_BEFORE, stateDUTBefore},  
  {STATE_DUT_DURING, stateDUTDuring},  
};

int numStates = 8;

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
  setRFPort(PORT_SHORT);
  reportRFPort(name_short);
  blink = BLINK_SHORT;
}

void stateOpenBefore(void) {

  state = STATE_OPEN_DURING;
  setRFPort(PORT_OPEN);
  reportRFPort(name_open);
  blink = BLINK_OPEN;
}

void stateLoadBefore(void) {

  state = STATE_LOAD_DURING;
  setRFPort(PORT_LOAD);
  reportRFPort(name_load);
  blink = BLINK_LOAD;
}

void stateDUTBefore(void) {

  state = STATE_DUT_DURING;
  setRFPort(PORT_DUT);
  reportRFPort(name_dut);
  blink = BLINK_DUT;
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

void stateDUTDuring(void) {

  state = STATE_DUT_DURING;
  // do nothing
}


void setRFPort(int port){

  /*
   *  Table 5. Truth Table
   * State V1 V2
   * RF1 on 0 0
   * RF2 on 1 0
   * RF3 on 0 1
   * RF4 on 1 1
   * 
  */

  switch(port) {
    case 1:
      digitalWrite(SWITCH_V1, LOW);
      digitalWrite(SWITCH_V2, LOW); 
      break;
    case 2:
      digitalWrite(SWITCH_V1, HIGH);
      digitalWrite(SWITCH_V2, LOW);
      break;
    case 3:
      digitalWrite(SWITCH_V1, LOW);
      digitalWrite(SWITCH_V2, HIGH);   
      break;
    case 4:  
      digitalWrite(SWITCH_V1, HIGH);  
      digitalWrite(SWITCH_V2, HIGH); 
      break;
    default:
      Serial.print("{\"report\":\"error\",\"is\":\"Port ");
      Serial.print(port);
      Serial.println(" is not known\"}");
              
  } 
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
  // SHORT: (10)(00)(00)(00)(delay)
  // OPEN:  (10)(10)(00)(00)(delay)
  // LOAD:  (10)(10)(10)(00)(delay)
  // DUT:   (10)(10)(10)(10)(delay) 
  // where delay is the same length again.

  switch (count%12){
      case 0:
        return true;
      case 2:
        switch (blink){
            case BLINK_OPEN:
            case BLINK_LOAD:
            case BLINK_DUT:
              return true;            
            default:
              return false;  
        }
      case 4:
        switch (blink){
            case BLINK_LOAD:
            case BLINK_DUT:
              return true;            
            default:
              return false;  
        }  
       case 6:
        switch (blink){
            case BLINK_DUT:
              return true;            
            default:
              return false;  
        }              
      default: 
        return false;
  }
  

  
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
    digitalWrite(LED_SWITCH,blinkState(count, blink));
    count++;
}

//===================================================================================
//====================== SETUP AND LOOP =============================================
//===================================================================================

void setup() {

  // pins for the RF switch control
  pinMode(SWITCH_V1, OUTPUT);
  pinMode(SWITCH_V2, OUTPUT);
  
  // pins for LED display
  pinMode(LED_SWITCH, OUTPUT);

  // Init timer ITimer1
  ITimer1.init();

  // timer
  ITimer1.attachInterruptInterval(TIMER_INTERVAL_MS, TimerHandler);
  
  Serial.setTimeout(50);
  Serial.begin(57600);
 
  while (! Serial); //wait for serial to start 

  count = 0; //initialise counter for use in blink code
 

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
        else if(strcmp(port, name_dut) == 0) {
          state = STATE_DUT_BEFORE;
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


 
