// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: query/query.proto

package org.zenoss.zing.proto.query;

public final class QueryOuterClass {
  private QueryOuterClass() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_query_Query_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_query_Query_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_query_Query_AdHocQuery_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_query_Query_AdHocQuery_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_query_Query_AdHocQuery_QueryParamsEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_query_Query_AdHocQuery_QueryParamsEntry_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_query_Query_AdHocQuery_MetadataEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_query_Query_AdHocQuery_MetadataEntry_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_query_Query_TypedQuery_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_query_Query_TypedQuery_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_query_Query_MetricQuery_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_query_Query_MetricQuery_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_query_Result_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_query_Result_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_query_Result_MetadataEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_query_Result_MetadataEntry_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_query_Result_MetricResult_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_query_Result_MetricResult_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n\021query/query.proto\022\005query\"\241\004\n\005Query\022\n\n\002" +
      "id\030\002 \001(\t\022/\n\014ad_hoc_query\030\003 \001(\0132\027.query.Q" +
      "uery.AdHocQueryH\000\022.\n\013typed_query\030\004 \001(\0132\027" +
      ".query.Query.TypedQueryH\000\0220\n\014metric_quer" +
      "y\030\005 \001(\0132\030.query.Query.MetricQueryH\000\032\200\002\n\n" +
      "AdHocQuery\022\024\n\014query_string\030\001 \001(\t\022>\n\014quer" +
      "y_params\030\002 \003(\0132(.query.Query.AdHocQuery." +
      "QueryParamsEntry\0227\n\010metadata\030\003 \003(\0132%.que" +
      "ry.Query.AdHocQuery.MetadataEntry\0322\n\020Que" +
      "ryParamsEntry\022\013\n\003key\030\001 \001(\t\022\r\n\005value\030\002 \001(",
      "\t:\0028\001\032/\n\rMetadataEntry\022\013\n\003key\030\001 \001(\t\022\r\n\005v" +
      "alue\030\002 \001(\t:\0028\001\032\032\n\nTypedQuery\022\014\n\004type\030\001 \001" +
      "(\t\032L\n\013MetricQuery\022\021\n\tmetric_id\030\001 \001(\t\022\r\n\005" +
      "start\030\002 \001(\003\022\013\n\003end\030\003 \001(\003\022\016\n\006latest\030\004 \001(\005" +
      "B\014\n\nquery_type\"\261\002\n\006Result\022\n\n\002id\030\001 \001(\t\022\017\n" +
      "\007queryId\030\002 \001(\t\022-\n\010metadata\030\003 \003(\0132\033.query" +
      ".Result.MetadataEntry\022\021\n\ttimestamp\030\004 \001(\003" +
      "\0223\n\rmetric_result\030\005 \001(\0132\032.query.Result.M" +
      "etricResultH\000\032/\n\rMetadataEntry\022\013\n\003key\030\001 " +
      "\001(\t\022\r\n\005value\030\002 \001(\t:\0028\001\032X\n\014MetricResult\022\023",
      "\n\013metric_name\030\001 \001(\t\022\021\n\tmetric_id\030\002 \001(\t\022\r" +
      "\n\005value\030\003 \001(\001\022\021\n\ttimestamp\030\004 \001(\003B\010\n\006resu" +
      "ltB\037\n\033org.zenoss.zing.proto.queryP\001b\006pro" +
      "to3"
    };
    com.google.protobuf.Descriptors.FileDescriptor.InternalDescriptorAssigner assigner =
        new com.google.protobuf.Descriptors.FileDescriptor.    InternalDescriptorAssigner() {
          public com.google.protobuf.ExtensionRegistry assignDescriptors(
              com.google.protobuf.Descriptors.FileDescriptor root) {
            descriptor = root;
            return null;
          }
        };
    com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
        }, assigner);
    internal_static_query_Query_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_query_Query_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_query_Query_descriptor,
        new java.lang.String[] { "Id", "AdHocQuery", "TypedQuery", "MetricQuery", "QueryType", });
    internal_static_query_Query_AdHocQuery_descriptor =
      internal_static_query_Query_descriptor.getNestedTypes().get(0);
    internal_static_query_Query_AdHocQuery_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_query_Query_AdHocQuery_descriptor,
        new java.lang.String[] { "QueryString", "QueryParams", "Metadata", });
    internal_static_query_Query_AdHocQuery_QueryParamsEntry_descriptor =
      internal_static_query_Query_AdHocQuery_descriptor.getNestedTypes().get(0);
    internal_static_query_Query_AdHocQuery_QueryParamsEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_query_Query_AdHocQuery_QueryParamsEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    internal_static_query_Query_AdHocQuery_MetadataEntry_descriptor =
      internal_static_query_Query_AdHocQuery_descriptor.getNestedTypes().get(1);
    internal_static_query_Query_AdHocQuery_MetadataEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_query_Query_AdHocQuery_MetadataEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    internal_static_query_Query_TypedQuery_descriptor =
      internal_static_query_Query_descriptor.getNestedTypes().get(1);
    internal_static_query_Query_TypedQuery_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_query_Query_TypedQuery_descriptor,
        new java.lang.String[] { "Type", });
    internal_static_query_Query_MetricQuery_descriptor =
      internal_static_query_Query_descriptor.getNestedTypes().get(2);
    internal_static_query_Query_MetricQuery_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_query_Query_MetricQuery_descriptor,
        new java.lang.String[] { "MetricId", "Start", "End", "Latest", });
    internal_static_query_Result_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_query_Result_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_query_Result_descriptor,
        new java.lang.String[] { "Id", "QueryId", "Metadata", "Timestamp", "MetricResult", "Result", });
    internal_static_query_Result_MetadataEntry_descriptor =
      internal_static_query_Result_descriptor.getNestedTypes().get(0);
    internal_static_query_Result_MetadataEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_query_Result_MetadataEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    internal_static_query_Result_MetricResult_descriptor =
      internal_static_query_Result_descriptor.getNestedTypes().get(1);
    internal_static_query_Result_MetricResult_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_query_Result_MetricResult_descriptor,
        new java.lang.String[] { "MetricName", "MetricId", "Value", "Timestamp", });
  }

  // @@protoc_insertion_point(outer_class_scope)
}
