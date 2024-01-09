<template>
  
 <div class='container-sm m-2 background-white border rounded'>

     <div v-if='!getCalibrated || !getVerified'>
         <h4 class='txt-primary txt-lg'>Please calibrate and verify the Pocket VNA before progressing</h4>
     </div>

    <div v-else>
        <!-- Drag and Drop elements -->
        <drag-and-drop-components id='dragdropcomponents' header="Device Under Test (DUT)" :display='duts' :syncPorts="getSyncPorts" @port1change='updatePort1' @port2change='updatePort2'/>
        
        <div class='d-flex flex-row justify-content-center'>

            <button id="measure" type='button' class="button-lg button-success" @click='measure' :disabled="port1 == '' && port2 == ''">Measure</button>

        </div>


        



        <transition name='fade'>
        <div v-if='getShowRequestModal' class="modal" id='modal-show' tabindex="-1">
            <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header" style="background-color: #ccc">
                <h5 class="modal-title">Requesting Measurement from Pocket VNA</h5>
                </div>
                <div class="modal-body">
                <div class='d-flex row align-items-center'>
                    <div class='col-2'>
                    <div class="spinner-border text-primary text-center" role="status">
                        <span class="visually-hidden">Loading...</span>
                    </div>
                    </div>
                    <div class='col-10'>
                    <p> Please wait for the pocket VNA to measure.</p>
                    <p>This could take around 30 seconds.</p>
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
</div>


</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import DragAndDropComponents from '@/components/DragAndDropComponents.vue';
import dayjs from "dayjs";

export default {
   name: 'DragAndDropMeasurement',
    props: ['calibrationState'],
    components:{
        DragAndDropComponents,
    },
    data () {
        return {
            duts: [
                {type: 'dut1', img: "dut1"},
                {type: 'dut2', img: "dut2"},
                {type: 'dut3', img: "dut3"},
                {type: 'dut4', img: "dut4"}
            ],

            port1: '',
            port2: '',

        }
    },
    computed:{
        ...mapGetters([
            'getCalibrated',
            'getVerified',
            'getShowRequestModal',
            'getSyncPorts'
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
            'setShowRequestModal',
        ]),
        measure(){
            this.rangeFreqRequest();
        }, 
        
      rangeFreqRequest(){
        //command structure: {"cmd":"crq","avg":1,"sparam":{"S11":true,"S12":false,"S21":false,"S22":false}}
        let params = {
            t:dayjs().unix(),
            avg:Number(this.calibrationState.avg),
            sparam:this.calibrationState.sparam,
            what: this.port1 != '' ? this.port1 : this.port2 
          }
          this.$store.dispatch('requestRangeAfterCal', params);
          this.$store.dispatch('setShowRequestModal', true);
      },
        updatePort1(connected){
            if(connected.type){
                this.port1 = connected.type;
            } else{
                this.port1 = '';
            }
        },
        updatePort2(connected){
            if(connected.type){
                this.port2 = connected.type;
            } else{
                this.port2 = '';
            }
            
        }


    }
}
</script>

<style>

</style>