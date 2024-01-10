<template>
  
 <div class='container-sm m-2 p-2 background-white border rounded'>
    <h4 v-if="sparams == 's11'"> 1 port Calibration Parameters </h4>
    <h4 v-else> 2 port Calibration Parameters </h4>

    <div v-if='sparams.length > 1' class='row mb-2'>
      <div class='col-12'>
        <div v-if='sparams.includes("s11")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s11check" checked disabled>
          <label class="form-check-label" for="s11check">S11</label>
        </div>

        <div v-if='sparams.includes("s12")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s12check" checked disabled>
          <label class="form-check-label" for="s12check">S12</label>
        </div>

        <div v-if='sparams.includes("s21")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s21check" checked disabled>
          <label class="form-check-label" for="s21check">S21</label>
        </div>

         <div v-if='sparams.includes("s22")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s22check" checked disabled>
          <label class="form-check-label" for="s22check">S22</label>
        </div>
        
      </div>
    </div>

    <div v-if="!isDisabled" class='row m-3' @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
      
        <div class='input-group'>
          <div class='col-md-2 pe-2'>
            <label for="sizeRange" class="txt-primary">Data points</label>
          </div>
          <div class='col-md-8 pe-2'>
            <input type="range" class="form-range" min="2" max="501" step="1" id="sizeRange" v-model='frequency_points' @change='setCalibrated(false); setVerified(false); setParametersSet(false)'>
          </div>
          <div class='col-md-2'>
            <label class='txt-primary'>{{frequency_points}}</label>
          </div>
        </div>
      
    </div>

    <div v-else class='row m-3'>
        
        <div class='input-group'>
        <div class='col-md-2 pe-2'>
            <label for="sizeRange" class="txt-grey">Data points</label>
        </div>
        <div class='col-md-8 pe-2'>
            <input type="range" class="form-range" min="2" max="501" step="1" id="sizeRange" :value='calibrationState.size' disabled>
        </div>
        <div class='col-md-2'>
            <label class='txt-grey'>{{ calibrationState.size }}</label>
        </div>
        </div>
    
    </div>

    
    <div v-if="!isDisabled" class='d-flex flex-column'>
        <!-- CALIBRATION COMMANDS-->
        <div class="input-group mb-2" @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
          <span class="input-group-text txt-primary col-sm-3">Start</span>
          <input type="number" :class="(parseFloat(frequency) >= 0 && parseFloat(frequency) < maxFrequency) ? 'form-control' : 'form-control is-invalid'" :min='minFrequency' :max='maxFrequency' placeholder="Start frequency" aria-label="freq" aria-describedby="basic-addon1" id="freq" v-model="frequency" @change='setCalibrated(false); setVerified(false); setParametersSet(false)' :disabled="isDisabled">
          <span class="input-group-text txt-primary" id="basic-addon1">MHz</span>
        </div>

        <div class="input-group" @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
          <span class="input-group-text txt-primary col-sm-3" id="basic-addon1">End</span>
          <input type="number" :class="(parseFloat(frequency_end) >= 0 && parseFloat(frequency_end) <= maxFrequency && parseFloat(frequency_end) > frequency) ? 'form-control' : 'form-control is-invalid'" :min='minFrequency' :max='maxFrequency' placeholder="End frequency" aria-label="freq" aria-describedby="basic-addon1" id="freq_end" v-model="frequency_end" @change='setCalibrated(false); setVerified(false); setParametersSet(false)' :disabled="isDisabled">     
          <span class="input-group-text txt-primary" id="basic-addon1">MHz</span>
        </div>
    </div>

    <div v-else class='d-flex flex-column'>
            <!-- CALIBRATION COMMANDS SHOWING BUT DISABLED-->
            <div class="input-group mb-2">
            <span class="input-group-text txt-grey col-sm-3">Start</span>
            <input type="number" class='form-control' placeholder="Start frequency" aria-label="freq" aria-describedby="basic-addon1" id="freq" :value="calibrationState.range.start/units" disabled>
            <span class="input-group-text txt-grey" id="basic-addon1">MHz</span>
            </div>

            <div class="input-group">
            <span class="input-group-text txt-grey col-sm-3" id="basic-addon1">End</span>
            <input type="number" class='form-control' aria-label="freq" aria-describedby="basic-addon1" id="freq_end" :value='calibrationState.range.end/units' disabled>     
            <span class="input-group-text txt-grey" id="basic-addon1">MHz</span>
            </div>
        </div>



    <div v-if="!isDisabled" class='d-flex flex-row justify-content-center'>
          <button id="set" type='button' class="button-lg button-success" @click='setCalibrationState'>Set Parameters</button>
    </div>


   
    <transition name='fade'>
        <div v-if='getShowParametersSetModal' class="modal" id='modal-show' tabindex="-1">
            <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header" style="background-color: #ccc">
                <h5 class="modal-title">Pocket VNA calibration parameters set</h5>
                </div>
                <div class="modal-body">
                <div class='d-flex row align-items-center'>
                    <div class='col-2' id="reveal-tick">
                        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" fill="green" class="bi bi-check-lg" viewBox="0 0 16 16">
                            <path d="M12.736 3.97a.733.733 0 0 1 1.047 0c.286.289.29.756.01 1.05L7.88 12.01a.733.733 0 0 1-1.065.02L3.217 8.384a.757.757 0 0 1 0-1.06.733.733 0 0 1 1.047 0l3.052 3.093 5.4-6.425z"/>
                        </svg>
                    </div>
                    <div class='col-10'>
                    <p> Calibration parameters have now been set.</p>
                    <p>Please move on to the next tab, 2) Calibration</p>
                    </div>
                </div>
                </div>
                <div class="modal-footer">
                <button type="button" class="btn btn-danger" @click="setShowParametersSetModal(false)">Close</button>
                </div>
            </div>
            </div>
        </div>
        </transition>

      

</div>


</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
   name: 'SetParameters',
    props: ['sparams', 'isDisabled', 'calibrationState'],
    components:{
        
    },
    data () {
        return {
            frequency_points: 20,
            frequency: 1.0,   //MHz
            frequency_end: 4.0, //MHz
            minFrequency: 0,
            maxFrequency: 4000.0, //MHz
            units: 1E6, //MHz
            avgCounts: 1,
            islog: false,
            s11: this.sparams.includes('s11') ? true:false,         
            s12: this.sparams.includes('s12') ? true:false,
            s21: this.sparams.includes('s21') ? true:false,
            s22: this.sparams.includes('s22') ? true:false,


        }
    },
    computed:{
        ...mapGetters([
            'getShowParametersSetModal'
        ]),
        
        
    },
    watch:{
        
    },
    created(){
        
    },
    mounted(){

    },
    methods:{
        ...mapActions([
            'setDraggable',
            'setCalibrated',
            'setVerified',
            'setParametersSet',
            'setShowParametersSetModal'
            
        ]),
        setCalibrationState(){
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
                range:{start: this.frequency*this.units, end: this.frequency_end*this.units},
                size: Number(this.frequency_points),
                avg:Number(this.avgCounts),
                islog: this.islog,
                sparam:{s11:this.s11,s12:this.s12,s21:this.s21,s22:this.s22}
            }

            this.$store.dispatch('setCalibrationState', params);
            this.$store.dispatch('sendCalibrationParameters', params);

        }
        


    }
}
</script>

<style>

</style>