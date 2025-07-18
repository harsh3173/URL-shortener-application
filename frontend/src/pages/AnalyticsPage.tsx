import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { urlApi } from '@/services/api';
import { Analytics, URLStats } from '@/types';
import LoadingSpinner from '@/components/LoadingSpinner';
import { 
  ArrowLeft, 
  MousePointer, 
  Users, 
  Calendar,
  Globe,
  Smartphone,
  Monitor,
  Tablet,
  ExternalLink
} from 'lucide-react';
import { format } from 'date-fns';
import { 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  LineChart,
  Line
} from 'recharts';
import toast from 'react-hot-toast';

const AnalyticsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [analytics, setAnalytics] = useState<Analytics[]>([]);
  const [stats, setStats] = useState<URLStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (id) {
      fetchAnalytics();
    }
  }, [id]);

  const fetchAnalytics = async () => {
    if (!id) return;
    
    setLoading(true);
    try {
      const response = await urlApi.getURLAnalytics(parseInt(id));
      setAnalytics(response.data!.analytics);
      setStats(response.data!.stats);
    } catch (error) {
      console.error('Error fetching analytics:', error);
      toast.error('Failed to fetch analytics');
    } finally {
      setLoading(false);
    }
  };

  const deviceData = analytics.reduce((acc, curr) => {
    const device = curr.device || 'Unknown';
    acc[device] = (acc[device] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const osData = analytics.reduce((acc, curr) => {
    const os = curr.os || 'Unknown';
    acc[os] = (acc[os] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const browserData = analytics.reduce((acc, curr) => {
    const browser = curr.browser || 'Unknown';
    acc[browser] = (acc[browser] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const dailyClicks = analytics.reduce((acc, curr) => {
    const date = format(new Date(curr.clicked_at), 'yyyy-MM-dd');
    acc[date] = (acc[date] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const chartData = Object.entries(dailyClicks).map(([date, clicks]) => ({
    date: format(new Date(date), 'MMM d'),
    clicks
  })).slice(-7);

  const pieData = Object.entries(deviceData).map(([device, count]) => ({
    name: device,
    value: count
  }));

  const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6'];

  const getDeviceIcon = (device: string) => {
    switch (device.toLowerCase()) {
      case 'mobile':
        return <Smartphone className="h-4 w-4" />;
      case 'tablet':
        return <Tablet className="h-4 w-4" />;
      case 'desktop':
        return <Monitor className="h-4 w-4" />;
      default:
        return <Globe className="h-4 w-4" />;
    }
  };

  const getBrowserIcon = (browser: string) => {
    switch (browser.toLowerCase()) {
      case 'chrome':
        return <Globe className="h-4 w-4" />;
      case 'firefox':
        return <Globe className="h-4 w-4" />;
      case 'safari':
        return <Globe className="h-4 w-4" />;
      default:
        return <Globe className="h-4 w-4" />;
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Analytics not found</h1>
          <Link to="/dashboard" className="btn btn-primary">
            Back to Dashboard
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <Link to="/dashboard" className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-4">
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Dashboard
        </Link>
        
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Analytics</h1>
            <p className="text-gray-600">Detailed statistics for your shortened URL</p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="card-content">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Clicks</p>
                <p className="text-2xl font-bold text-gray-900">{stats.total_clicks}</p>
              </div>
              <MousePointer className="h-8 w-8 text-blue-600" />
            </div>
          </div>
        </div>
        
        <div className="card">
          <div className="card-content">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Unique Visitors</p>
                <p className="text-2xl font-bold text-gray-900">{stats.unique_clicks}</p>
              </div>
              <Users className="h-8 w-8 text-green-600" />
            </div>
          </div>
        </div>
        
        <div className="card">
          <div className="card-content">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Last Clicked</p>
                <p className="text-2xl font-bold text-gray-900">
                  {stats.last_clicked ? format(new Date(stats.last_clicked), 'MMM d') : 'Never'}
                </p>
              </div>
              <Calendar className="h-8 w-8 text-purple-600" />
            </div>
          </div>
        </div>
        
        <div className="card">
          <div className="card-content">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Click Rate</p>
                <p className="text-2xl font-bold text-gray-900">
                  {stats.unique_clicks > 0 ? (stats.total_clicks / stats.unique_clicks).toFixed(1) : '0'}
                </p>
              </div>
              <ExternalLink className="h-8 w-8 text-orange-600" />
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="card">
          <div className="card-header">
            <h2 className="text-xl font-semibold text-gray-900">Clicks Over Time</h2>
          </div>
          <div className="card-content">
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip />
                <Line type="monotone" dataKey="clicks" stroke="#3b82f6" strokeWidth={2} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="card">
          <div className="card-header">
            <h2 className="text-xl font-semibold text-gray-900">Device Types</h2>
          </div>
          <div className="card-content">
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {pieData.map((_entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="card">
          <div className="card-header">
            <h2 className="text-xl font-semibold text-gray-900">Top Devices</h2>
          </div>
          <div className="card-content">
            <div className="space-y-3">
              {Object.entries(deviceData)
                .sort(([,a], [,b]) => b - a)
                .slice(0, 5)
                .map(([device, count]) => (
                  <div key={device} className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      {getDeviceIcon(device)}
                      <span className="text-sm font-medium text-gray-900">{device}</span>
                    </div>
                    <span className="text-sm text-gray-600">{count}</span>
                  </div>
                ))}
            </div>
          </div>
        </div>

        <div className="card">
          <div className="card-header">
            <h2 className="text-xl font-semibold text-gray-900">Top Browsers</h2>
          </div>
          <div className="card-content">
            <div className="space-y-3">
              {Object.entries(browserData)
                .sort(([,a], [,b]) => b - a)
                .slice(0, 5)
                .map(([browser, count]) => (
                  <div key={browser} className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      {getBrowserIcon(browser)}
                      <span className="text-sm font-medium text-gray-900">{browser}</span>
                    </div>
                    <span className="text-sm text-gray-600">{count}</span>
                  </div>
                ))}
            </div>
          </div>
        </div>

        <div className="card">
          <div className="card-header">
            <h2 className="text-xl font-semibold text-gray-900">Top Operating Systems</h2>
          </div>
          <div className="card-content">
            <div className="space-y-3">
              {Object.entries(osData)
                .sort(([,a], [,b]) => b - a)
                .slice(0, 5)
                .map(([os, count]) => (
                  <div key={os} className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <Globe className="h-4 w-4" />
                      <span className="text-sm font-medium text-gray-900">{os}</span>
                    </div>
                    <span className="text-sm text-gray-600">{count}</span>
                  </div>
                ))}
            </div>
          </div>
        </div>
      </div>

      {analytics.length > 0 && (
        <div className="card mt-6">
          <div className="card-header">
            <h2 className="text-xl font-semibold text-gray-900">Recent Activity</h2>
          </div>
          <div className="card-content">
            <div className="space-y-2">
              {analytics.slice(0, 10).map((click, index) => (
                <div key={index} className="flex items-center justify-between py-2 border-b border-gray-100 last:border-b-0">
                  <div className="flex items-center space-x-3">
                    {getDeviceIcon(click.device || 'Unknown')}
                    <div>
                      <div className="text-sm font-medium text-gray-900">
                        {click.device || 'Unknown'} • {click.os || 'Unknown'} • {click.browser || 'Unknown'}
                      </div>
                      <div className="text-xs text-gray-500">
                        {click.ip_address} • {click.country || 'Unknown'}
                      </div>
                    </div>
                  </div>
                  <div className="text-sm text-gray-600">
                    {format(new Date(click.clicked_at), 'MMM d, HH:mm')}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AnalyticsPage;