
<template>
    <div>
        <h4> 2 port Calibration and Measurement </h4>

    <div v-if='singleFrequencyAllowed' class='row mb-2'>
      <div class='col-12'>
        <div v-if='singleFrequencyAllowed' class="form-check form-check-inline">
          <input class="form-check-input" type="radio" name="singleRadio" value='single' id="singleRadio" v-model='mode' :disabled='!rangeFrequencyAllowed'>
          <label class="form-check-label" for="singleRadio">Single Frequency</label>
        </div>
        <div v-if='rangeFrequencyAllowed' class="form-check form-check-inline">
          <input class="form-check-input" type="radio" name="rangeRadio" value='range' id="rangeRadio" v-model='mode' :disabled='!singleFrequencyAllowed'>
          <label class="form-check-label" for="rangeRadio">Range of frequencies</label>
        </div>
      </div>
    </div>


    <div class='row mb-2'>
      <div class='col-12 d-flex'>
        <div class="form-check form-check-inline">
          <input class="form-check-input" type="radio" name="hardware-radio1" value='filter' id="hardware-radio1" v-model='hardware' @change="updateHardwareSelected('filter')">
          <label class="form-check-label" for="hardware-radio1">Filter</label>
        </div>

        <div class="form-check form-check-inline">
          <input class="form-check-input" type="radio" name="hardware-radio2" value='coupler' id="hardware-radio2" v-model='hardware' @change="updateHardwareSelected('coupler')">
          <label class="form-check-label" for="hardware-radio2">Coupler</label>
        </div>
      </div>
    </div>

    <!-- <div v-if='sparams.length > 1' class='row mb-2'>
      <div class='col-12'>
        <div v-if='sparams.includes("s11")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s11check" v-model='s11' disabled >
          <label class="form-check-label" for="s11check">S11</label>
        </div>

        <div v-if='sparams.includes("s12")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s12check" v-model='s12' disabled>
          <label class="form-check-label" for="s12check">S12</label>
        </div>

        <div v-if='sparams.includes("s21")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s21check" v-model='s21' disabled>
          <label class="form-check-label" for="s21check">S21</label>
        </div>

         <div v-if='sparams.includes("s22")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s22check" v-model='s22' disabled>
          <label class="form-check-label" for="s22check">S22</label>
        </div>
        
      </div>
    </div> -->

    <div v-if='averagingAllowed' class='row mb-2' @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
      <div class='col-12'>
        <div class='input-group'>
          <div class='col-3'>
            <span for="averageRange" class="input-group-text me-4">Average</span>
          </div>
          <div class='col-6'>
            <input type="range" class="form-range" min="1" max="10" step="1" id="averageRange" v-model='avgCounts'>
          </div>
          <div class='col-2'>
            <span class='input-group-text ms-4'>{{avgCounts}}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if='mode == "range"' class='row p-3' @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
      <div class='col-12'>
        <div class='input-group'>
          <div class='col-lg-3'>
            <span for="sizeRange" class="input-group-text me-4 text-wrap">Frequency points</span>
          </div>
          <div class='col-lg-7'>
            <input type="range" class="form-range" :min="min_data_points" :max="max_data_points" step="1" id="sizeRange" v-model='size' @change='setCalibrated(false)'>
          </div>
          <div class='col-lg-2'>
            <span class='input-group-text ms-4'>{{size}}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- <div v-if='mode == "range"' class='row mb-2'>
      <div class="form-check">
        <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="isLogcheck" v-model='islog' :disabled='getCalibrated'>
        <label class="form-check-label" for="flexCheckDefault">
            Logarithmic
        </label>
      </div>
    </div> -->

    
    <div class='row p-3'>
      <div class='col-lg-10'>
        <!-- CALIBRATION COMMANDS-->
        <div v-if='mode == "single"' class="input-group" @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
          <span class="input-group-text" id="basic-addon1">Frequency</span>
          <input type="number" :class="(parseFloat(frequency) >= 0 && parseFloat(frequency) < maxFrequency) ? 'form-control' : 'form-control is-invalid'" :min='minFrequency' :max='maxFrequency' placeholder="frequency" aria-label="freq" aria-describedby="basic-addon1" id="freq" v-model="frequency">
          <span class="input-group-text" id="basic-addon1">MHz</span>
          <button id="request" type='button' class="btn btn-success btn-lg" @click="singleFreqCommand">Request</button>
        </div>

        <div v-else @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
            <div class="input-group">
                <span class="input-group-text" id="span-start-1">Start</span>
                <input type="number" :class="(parseFloat(frequency) >= 0 && parseFloat(frequency) < maxFrequency) ? 'form-control' : 'form-control is-invalid'" :min='minFrequency' :max='maxFrequency' placeholder="Start frequency" aria-label="freq" aria-describedby="span-start-1" id="freq" v-model="frequency" @change="setCalibrated(false)">
                <span class="input-group-text" id="span-end-1">MHz</span>
            </div>

            <div class="input-group">
                <span class="input-group-text" id="span-start-2">End&nbsp;</span>
                <input type="number" :class="(parseFloat(frequency_end) >= 0 && parseFloat(frequency_end) <= maxFrequency && parseFloat(frequency_end) > frequency) ? 'form-control' : 'form-control is-invalid'" :min='minFrequency' :max='maxFrequency' placeholder="End frequency" aria-label="freq" aria-describedby="span-start-2" id="freq_end" v-model="frequency_end" @change="setCalibrated(false)">     
                <span class="input-group-text" id="span-end-2">MHz</span>
            </div>
          

          
          
        </div>
      </div>

      <div class='d-flex col-lg-2'>
          <button id="request_calibration" type='button' class="btn btn-success btn-lg" @click="rangeFreqCalibration" :disabled='getSessionExpired'>Calibrate</button>
      </div>
    </div>

    


    <div class='row p-3'>
      <div class='col-lg-10'>

         <!-- MEASUREMENT COMMANDS-->
        <div v-if='mode == "single"' class="input-group mb-2" @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
          <span class="input-group-text" id="basic-addon1">Frequency</span>
          <input type="number" :class="(parseFloat(frequency) >= 0 && parseFloat(frequency) < maxFrequency) ? 'form-control' : 'form-control is-invalid'" :min='minFrequency' :max='maxFrequency' placeholder="frequency" aria-label="freq" aria-describedby="basic-addon1" id="freq" v-model="frequency">
          <span class="input-group-text" id="basic-addon1">MHz</span>
          <button id="request" type='button' class="btn btn-success btn-lg" @click="singleFreqCommand">Request</button>
        </div>

        <div v-else @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
            <div class="input-group">
                <span class="input-group-text" id="span-start-measure-1">Start</span>
                <input type="number" :class="(parseFloat(frequency) >= 0 && parseFloat(frequency) < maxFrequency) ? 'form-control' : 'form-control is-invalid'" :min='minFrequency' :max='maxFrequency' placeholder="Start frequency" aria-label="freq" aria-describedby="span-start-measure-1" id="freq" :value="frequency" disabled>
                <span class="input-group-text" id="span-end-measure-1">MHz</span>
            </div>

            <div class="input-group mb-2">
                <span class="input-group-text" id="span-start-measure-2">End&nbsp;</span>
                <input type="number" :class="(parseFloat(frequency_end) >= 0 && parseFloat(frequency_end) <= maxFrequency && parseFloat(frequency_end) > frequency) ? 'form-control' : 'form-control is-invalid'" :min='minFrequency' :max='maxFrequency' placeholder="End frequency" aria-label="freq" aria-describedby="span-start-measure-2" id="freq_end" :value="frequency_end" disabled>     
                <span class="input-group-text" id="span-end-measure-2">MHz</span>
            </div>
          
          
        </div>

        <!-- WHAT IS BEING MEASURED-->
        <!-- <div class='row mb-2'>
          <div class='col-12 d-flex justify-content-end'>
            <div class="form-check form-check-inline">
              <input class="form-check-input" type="radio" name="dutRadio1" value='dut1' id="dutRadio1" v-model='what'>
              <label class="form-check-label" for="dutRadio1">DUT1</label>
            </div>
            <div class="form-check form-check-inline">
              <input class="form-check-input" type="radio" name="dutRadio2" value='dut2' id="dutRadio2" v-model='what'>
              <label class="form-check-label" for="dutRadio2">DUT2</label>
            </div>
            <div class="form-check form-check-inline">
              <input class="form-check-input" type="radio" name="dutRadio3" value='dut3' id="dutRadio3" v-model='what'>
              <label class="form-check-label" for="dutRadio3">DUT3</label>
            </div>
            <div class="form-check form-check-inline">
              <input class="form-check-input" type="radio" name="dutRadio4" value='dut4' id="dutRadio4" v-model='what'>
              <label class="form-check-label" for="dutRadio4">DUT4</label>
            </div>
            <div class="form-check form-check-inline">
              <input class="form-check-input" type="radio" name="shortRadio" value='short' id="shortRadio" v-model='what'>
              <label class="form-check-label" for="shortRadio">Short</label>
            </div>
            <div class="form-check form-check-inline">
              <input class="form-check-input" type="radio" name="openRadio" value='open' id="openRadio" v-model='what'>
              <label class="form-check-label" for="openRadio">Open</label>
            </div>
            <div class="form-check form-check-inline">
              <input class="form-check-input" type="radio" name="loadRadio" value='load' id="loadRadio" v-model='what'>
              <label class="form-check-label" for="loadRadio">Load</label>
            </div>
            <div class="form-check form-check-inline">
              <input class="form-check-input" type="radio" name="thruRadio" value='thru' id="thruRadio" v-model='what'>
              <label class="form-check-label" for="thruRadio">Thru</label>
            </div>
          </div>
        </div> -->

        <div class="row">
          <div class='input-group'>
              <span class='input-group-text' for="graph">Select measurement: </span>
              <select class='form-select' name="measurement" id="measurement" v-model="what" @change="updateDUTSelected(what)">
                  <option v-if="hardware == 'filter'" value="dut1">Port 1 and 2</option>
                  <option v-if="hardware == 'coupler'" value="dut2">Port 1 and 4</option>
                  <option v-if="hardware == 'coupler'" value="dut3">Port 1 and 2</option>
                  <option v-if="hardware == 'coupler'" value="dut4">Port 1 and 3</option>
                  <!-- <option value="short">Short</option>
                  <option value="open">Open</option>
                  <option value="load">Load</option>
                  <option value="thru">Thru</option> -->
                  
              </select> 
          </div>
        </div>
       

    </div>

      <div class='d-flex col-lg-2'>
        <button id="request_results" type='button' class="btn btn-success btn-lg" @click="rangeFreqRequest" :disabled='!getCalibrated || getSessionExpired'>Measure</button>
      </div>

  </div>


    <transition name='fade'>
      <div v-if='getShowCalibrationModal' class="modal" id='modal-show' tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header" style="background-color: #ccc">
              <h5 class="modal-title">Calibrating Pocket VNA</h5>
              <!-- <button type="button" class="btn btn-close" aria-label="Close" @click='showLoadDataModal = false'>
                
              </button> -->
            </div>
            <div class="modal-body">
              <div class='d-flex row align-items-center'>
                <div class='col-2'>
                  <div class="spinner-border text-primary text-center" role="status">
                    <span class="visually-hidden">Loading...</span>
                  </div>
                </div>
                <div class='col-10'>
                  <p> Please wait for the pocket VNA to calibrate.</p>
                  <p>This could take a couple of minutes.</p>
                  <p>If it takes much longer than this, please click Close, then refresh the page and try again.</p>
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-danger" @click="setShowCalibrationModal(false)">Close</button>
            </div>
          </div>
        </div>
      </div>
      </transition>

      <transition name='fade'>
      <div v-if='getShowRequestModal' class="modal" id='modal-show' tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header" style="background-color: #ccc">
              <h5 class="modal-title">Requesting measurements from Pocket VNA</h5>
              <!-- <button type="button" class="btn btn-close" aria-label="Close" @click='showLoadDataModal = false'>
                
              </button> -->
            </div>
            <div class="modal-body">
              <div class='d-flex row align-items-center'>
                <div class='col-2'>
                  <div class="spinner-border text-primary text-center" role="status">
                    <span class="visually-hidden">Loading...</span>
                  </div>
                </div>
                <div class='col-10'>
                  <p> Please wait for the pocket VNA to send data.</p>
                  <p>This could take a couple of minutes.</p>
                  <p>If it takes much longer than this, please click Close, then refresh the page and try again.</p>
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-danger" @click="setShowRequestModal(false)">Close</button>
            </div>
          </div>
        </div>
      </div>
      </transition>

</div>
</template>

<script>
import dayjs from "dayjs";
import { mapActions, mapGetters } from 'vuex';

export default {

  name: 'SendCommands',
  props:['singleFrequencyAllowed', 'rangeFrequencyAllowed', 'sparams', 'averagingAllowed'],
  data () {
    return {
      hardware: 'filter', // or 'coupler'
        mode: 'range',  // or 'range'
        frequency: 1.0,   //MHz
        frequency_end: 4.0, //MHz
        min_data_points: 2,
        max_data_points: 201,
        minFrequency: 0,
        maxFrequency: 4000.0, //MHz
        units: 1E6,
        s11: true,
        s12: true,
        s21: true,
        s22: true,
        avgCounts: 1,
        size: 20,
        islog: false,
        what: 'dut1', //'short', 'load', 'open', 'thru', 'dut1', 'dut2', 'dut3', 'dut4'

        //showCancel: false,  //should the cancel button show on the modals

    }
  },
  components: {
    
  },
  computed:{
      ...mapGetters([
        'getCalibrated',
        'getShowCalibrationModal',
        'getShowRequestModal',
        'getSessionExpired'
      ])
  },
  watch:{
    getSessionExpired(expired){
        if(expired){
            this.setShowCalibrationModal(false);
            this.setShowRequestModal(false);
        }
      }
  },
  created(){
      if(this.singleFrequencyAllowed && !this.rangeFrequencyAllowed){
        this.mode = 'single';
      } else if(!this.singleFrequencyAllowed && this.rangeFrequencyAllowed){
        this.mode = 'range';
      } else{
        this.mode = 'single';
      }
  },
  mounted(){
      

  },
  methods: {
    ...mapActions([
      'setDraggable',
      'setCalibrated',
      'setShowCalibrationModal',
      'setShowRequestModal',
      'setDUT'
    ]),
      singleFreqCommand(){
          //command structure: {"id":"945102d5-94e4-448e-bbbf-48384c662711","t":1634664795,"cmd":"sq","freq":100000,"avg":1,"sparam":{"S11":true,"S12":false,"S21":true,"S22":false}}
          if(this.frequency < this.minFrequency){
            this.frequency = this.minFrequency;
          } else if(this.frequency > this.maxFrequency){
            this.frequency = this.maxFrequency;
          }
          
          let params = {
            id:'',
            t:dayjs().unix(),
            freq:this.frequency*this.units,
            avg:this.avgCounts,
            sparam:{s11:this.s11,s12:this.s12,s21:this.s21,s22:this.s22}
          }
          this.$store.dispatch('requestSingle', params);
      },
      rangeFreqCommand(){
        //command structure: {"cmd":"rq","range":{"Start":100000,"End":4000000},"size":2,"isLog":true,"avg":1,"sparam":{"S11":true,"S12":false,"S21":true,"S22":false}}
        if(this.frequency < this.minFrequency){
            this.frequency = this.minFrequency;
          } else if(this.frequency > this.maxFrequency){
            this.frequency = this.maxFrequency - 1;
          }

          if(this.frequency_end < this.frequency){
            this.frequency_end = this.frequency + 1;
          } else if(this.frequency_end > this.maxFrequency){
            this.frequency_end = this.maxFrequency;
          }
        
        let params = {
            id:'',
            t:dayjs().unix(),
            range:{start: this.frequency*this.units, end: this.frequency_end*this.units},
            size: Number(this.size),
            islog: this.islog,
            avg:Number(this.avgCounts),
            sparam:{s11:this.s11,s12:this.s12,s21:this.s21,s22:this.s22}
          }

          this.$store.dispatch('requestRange', params);
      },
      rangeFreqCalibration(){
        //command structure: {"cmd":"rc","range":{"Start":100000,"End":4000000},"size":2,"isLog":true,"avg":1,"sparam":{"S11":true,"S12":false,"S21":true,"S22":false}}
        if(this.frequency < this.minFrequency){
            this.frequency = this.minFrequency;
          } else if(this.frequency > this.maxFrequency){
            this.frequency = this.maxFrequency - 1;
          }

          if(this.frequency_end < this.frequency){
            this.frequency_end = this.frequency + 1;
          } else if(this.frequency_end > this.maxFrequency){
            this.frequency_end = this.maxFrequency;
          }

        
        let params = {
            id:'',
            t:dayjs().unix(),
            range:{start: this.frequency*this.units, end: this.frequency_end*this.units},
            size: Number(this.size),
            islog: this.islog,
            avg:Number(this.avgCounts),
            sparam:{s11:this.s11,s12:this.s12,s21:this.s21,s22:this.s22}
          }

          this.$store.dispatch('requestCalibration', params);
          this.$store.dispatch('setShowCalibrationModal', true);
          
      },
      rangeFreqRequest(){
        //command structure: {"cmd":"crq","avg":1,"sparam":{"S11":true,"S12":false,"S21":false,"S22":false}}
        //only works for s11 parameter alone
        let params = {
            avg:Number(this.avgCounts),
            sparam:{s11:this.s11,s12:this.s12,s21:this.s21,s22:this.s22},
            what: this.what
          }

          this.$store.dispatch('requestRangeAfterCal', params);
          this.$store.dispatch('setShowRequestModal', true);
      },
      toggleCalibrationMessage(){
        this.showCalibrationMessage = !this.showCalibrationMessage;
      },
      updateDUTSelected(dut){
        if(dut == 'dut1' || dut == 'dut2' || dut == 'dut3' || dut == 'dut4'){
          this.setDUT(dut);
        } 
      },
      updateHardwareSelected(hardware){
        if(hardware == 'filter'){
          this.what = 'dut1';
          this.setDUT('dut1');
        } else{
          this.what = 'dut2';
          this.setDUT('dut2');
        }
      }
     
      
  }
}
</script>

<style scoped>
input[type=number]{
  min-width: 200px;
}

</style>