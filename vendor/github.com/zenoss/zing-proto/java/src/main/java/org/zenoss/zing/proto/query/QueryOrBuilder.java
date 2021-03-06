// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: query/query.proto

package org.zenoss.zing.proto.query;

public interface QueryOrBuilder extends
    // @@protoc_insertion_point(interface_extends:query.Query)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * A unique ID of this query (Will be included in the results)
   * </pre>
   *
   * <code>string id = 2;</code>
   */
  java.lang.String getId();
  /**
   * <pre>
   * A unique ID of this query (Will be included in the results)
   * </pre>
   *
   * <code>string id = 2;</code>
   */
  com.google.protobuf.ByteString
      getIdBytes();

  /**
   * <code>.query.Query.AdHocQuery ad_hoc_query = 3;</code>
   */
  boolean hasAdHocQuery();
  /**
   * <code>.query.Query.AdHocQuery ad_hoc_query = 3;</code>
   */
  org.zenoss.zing.proto.query.Query.AdHocQuery getAdHocQuery();
  /**
   * <code>.query.Query.AdHocQuery ad_hoc_query = 3;</code>
   */
  org.zenoss.zing.proto.query.Query.AdHocQueryOrBuilder getAdHocQueryOrBuilder();

  /**
   * <code>.query.Query.TypedQuery typed_query = 4;</code>
   */
  boolean hasTypedQuery();
  /**
   * <code>.query.Query.TypedQuery typed_query = 4;</code>
   */
  org.zenoss.zing.proto.query.Query.TypedQuery getTypedQuery();
  /**
   * <code>.query.Query.TypedQuery typed_query = 4;</code>
   */
  org.zenoss.zing.proto.query.Query.TypedQueryOrBuilder getTypedQueryOrBuilder();

  /**
   * <code>.query.Query.MetricQuery metric_query = 5;</code>
   */
  boolean hasMetricQuery();
  /**
   * <code>.query.Query.MetricQuery metric_query = 5;</code>
   */
  org.zenoss.zing.proto.query.Query.MetricQuery getMetricQuery();
  /**
   * <code>.query.Query.MetricQuery metric_query = 5;</code>
   */
  org.zenoss.zing.proto.query.Query.MetricQueryOrBuilder getMetricQueryOrBuilder();

  public org.zenoss.zing.proto.query.Query.QueryTypeCase getQueryTypeCase();
}
