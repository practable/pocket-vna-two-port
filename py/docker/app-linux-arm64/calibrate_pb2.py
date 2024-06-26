# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: calibrate.proto
# Protobuf Python Version: 4.25.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0f\x63\x61librate.proto\x12\x02pb\"J\n\x18\x43\x61librateOnePortResponse\x12\x11\n\tfrequency\x18\x01 \x03(\x01\x12\x1b\n\x06result\x18\x02 \x03(\x0b\x32\x0b.pb.Complex\"J\n\x18\x43\x61librateTwoPortResponse\x12\x11\n\tfrequency\x18\x01 \x03(\x01\x12\x1b\n\x06result\x18\x02 \x01(\x0b\x32\x0b.pb.SParams\"\xb3\x01\n\x17\x43\x61librateOnePortRequest\x12\x11\n\tfrequency\x18\x01 \x03(\x01\x12\x1a\n\x05short\x18\x02 \x03(\x0b\x32\x0b.pb.Complex\x12\x19\n\x04open\x18\x03 \x03(\x0b\x32\x0b.pb.Complex\x12\x19\n\x04load\x18\x04 \x03(\x0b\x32\x0b.pb.Complex\x12\x19\n\x04thru\x18\x05 \x03(\x0b\x32\x0b.pb.Complex\x12\x18\n\x03\x64ut\x18\x06 \x03(\x0b\x32\x0b.pb.Complex\"\xb3\x01\n\x17\x43\x61librateTwoPortRequest\x12\x11\n\tfrequency\x18\x01 \x03(\x01\x12\x1a\n\x05short\x18\x02 \x01(\x0b\x32\x0b.pb.SParams\x12\x19\n\x04open\x18\x03 \x01(\x0b\x32\x0b.pb.SParams\x12\x19\n\x04load\x18\x04 \x01(\x0b\x32\x0b.pb.SParams\x12\x19\n\x04thru\x18\x05 \x01(\x0b\x32\x0b.pb.SParams\x12\x18\n\x03\x64ut\x18\x06 \x01(\x0b\x32\x0b.pb.SParams\"q\n\x07SParams\x12\x18\n\x03s11\x18\x01 \x03(\x0b\x32\x0b.pb.Complex\x12\x18\n\x03s12\x18\x02 \x03(\x0b\x32\x0b.pb.Complex\x12\x18\n\x03s21\x18\x03 \x03(\x0b\x32\x0b.pb.Complex\x12\x18\n\x03s22\x18\x04 \x03(\x0b\x32\x0b.pb.Complex\"%\n\x07\x43omplex\x12\x0c\n\x04imag\x18\x01 \x01(\x01\x12\x0c\n\x04real\x18\x02 \x01(\x01\x32\xad\x01\n\tCalibrate\x12O\n\x10\x43\x61librateOnePort\x12\x1b.pb.CalibrateOnePortRequest\x1a\x1c.pb.CalibrateOnePortResponse\"\x00\x12O\n\x10\x43\x61librateTwoPort\x12\x1b.pb.CalibrateTwoPortRequest\x1a\x1c.pb.CalibrateTwoPortResponse\"\x00\x42\x31Z/github.com/practable/pocket-vna-two-port/pkg/pbb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'calibrate_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  _globals['DESCRIPTOR']._options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z/github.com/practable/pocket-vna-two-port/pkg/pb'
  _globals['_CALIBRATEONEPORTRESPONSE']._serialized_start=23
  _globals['_CALIBRATEONEPORTRESPONSE']._serialized_end=97
  _globals['_CALIBRATETWOPORTRESPONSE']._serialized_start=99
  _globals['_CALIBRATETWOPORTRESPONSE']._serialized_end=173
  _globals['_CALIBRATEONEPORTREQUEST']._serialized_start=176
  _globals['_CALIBRATEONEPORTREQUEST']._serialized_end=355
  _globals['_CALIBRATETWOPORTREQUEST']._serialized_start=358
  _globals['_CALIBRATETWOPORTREQUEST']._serialized_end=537
  _globals['_SPARAMS']._serialized_start=539
  _globals['_SPARAMS']._serialized_end=652
  _globals['_COMPLEX']._serialized_start=654
  _globals['_COMPLEX']._serialized_end=691
  _globals['_CALIBRATE']._serialized_start=694
  _globals['_CALIBRATE']._serialized_end=867
# @@protoc_insertion_point(module_scope)
